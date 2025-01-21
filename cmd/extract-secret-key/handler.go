package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/tidwall/gjson"
)

//go:embed templates/*
var templatesFS embed.FS

type TemplateData struct {
	Error  string  `json:"error"`
	Raw    string  `json:"raw"`
	Things []Thing `json:"things"`
}

type Thing struct {
	ThingID   string `json:"thingID"`
	ThingName string `json:"thingName"`
	SecretKey string `json:"secretKey"`
}

// renderTemplate renders a template.
func renderTemplate(w http.ResponseWriter, childTemplate string, data any) {
	childTemplate = "templates/" + childTemplate
	tmpl, err := template.ParseFS(templatesFS, "templates/base.gohtml", childTemplate)
	if err != nil {
		log.Printf("error parsing template %s: %v\n", childTemplate, err)
		http.Error(w, "error parsing templates: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base.gohtml", data); err != nil {
		log.Printf("error executing template %s: %v\n", childTemplate, err)
		http.Error(w, "error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// renderError renders an error.
func renderError(w http.ResponseWriter, err string) {
	renderTemplate(w, "error.gohtml", TemplateData{
		Error: err,
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.gohtml", nil)
}

const (
	awsRegion    = "us-east-1"
	poolID       = "us-east-1_s6T1EfUgP"
	clientID     = "2ut21oh12e4gb6t5f8g87ls0t3"
	clientSecret = "dtof6slltmbet3gic9o1cgirur7ietd4mqtklt9r0ld5h902la2"
	awsUserAgent = "aws-sdk-android/2.22.6 Linux/5.4.226 Dalvik/2.1.0/0 pt_BR"
)

func accountInfoHandler() func(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Transport: &userAgentTransport{
			UserAgent: awsUserAgent,
			Base:      http.DefaultTransport,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx      = r.Context()
			username = strings.TrimSpace(r.FormValue("username"))
			password = strings.TrimSpace(r.FormValue("password"))
		)
		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			renderError(w, "username or password is empty")
			return
		}

		// configure cognito srp
		csrp, err := cognitosrp.NewCognitoSRP(username, password, poolID, clientID, aws.String(clientSecret))
		if err != nil {
			log.Printf("error creating SRP: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			renderError(w, fmt.Sprintf("error creating SRP: %v", err))
			return
		}

		// configure cognito identity provider
		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithHTTPClient(httpClient),
			config.WithRegion(awsRegion),
			config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		)
		if err != nil {
			log.Printf("error loading AWS config: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			renderError(w, fmt.Sprintf("error loading AWS config: %v", err))
			return
		}

		svc := cip.NewFromConfig(cfg)

		// initialize the auth
		resp, err := svc.InitiateAuth(ctx, &cip.InitiateAuthInput{
			AuthFlow:       types.AuthFlowTypeUserSrpAuth,
			ClientId:       aws.String(csrp.GetClientId()),
			AuthParameters: csrp.GetAuthParams(),
		})
		if err != nil {
			log.Printf("error calling InitiateAuth: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			renderError(w, fmt.Sprintf("error calling InitiateAuth: %v", err))
			return
		}

		// we only support PASSWORD_VERIFIER challenge
		if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
			challengeResponses, err := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())
			if err != nil {
				log.Printf("error calling PasswordVerifierChallenge: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				renderError(w, fmt.Sprintf("error calling PasswordVerifierChallenge: %v", err))
				return
			}

			resp, err := svc.RespondToAuthChallenge(ctx, &cip.RespondToAuthChallengeInput{
				ChallengeName:      types.ChallengeNameTypePasswordVerifier,
				ChallengeResponses: challengeResponses,
				ClientId:           aws.String(csrp.GetClientId()),
			})
			if err != nil {
				log.Printf("error calling RespondToAuthChallenge: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				renderError(w, fmt.Sprintf("error calling RespondToAuthChallenge: %v", err))
				return
			}

			if resp.AuthenticationResult.AccessToken == nil || resp.AuthenticationResult.IdToken == nil {
				log.Printf("accessToken or idToken is nil\n")
				w.WriteHeader(http.StatusInternalServerError)
				renderError(w, "accessToken or idToken is nil")
				return
			}

			// get the user
			user, err := svc.GetUser(ctx, &cip.GetUserInput{
				AccessToken: resp.AuthenticationResult.AccessToken,
			})
			if err != nil {
				log.Printf("error calling GetUser: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				renderError(w, fmt.Sprintf("error calling GetUser: %v", err))
				return
			}

			// got nil username
			if user.Username == nil {
				log.Printf("user.Username is nil\n")
				w.WriteHeader(http.StatusInternalServerError)
				renderError(w, "user.Username is nil")
				return
			}

			// call the managething endpoint
			raw, err := managething(ctx, *user.Username, *resp.AuthenticationResult.AccessToken, *resp.AuthenticationResult.IdToken)
			if err != nil {
				// ensure we don't log PII data
				if !strings.HasPrefix(err.Error(), "invalidResponse") {
					log.Printf("error calling managething: %v\n", err)
				}
				w.WriteHeader(http.StatusInternalServerError)
				renderError(w, fmt.Sprintf("error calling managething: %v", err))
				return
			}

			// parse the things
			var things []Thing
			for _, t := range gjson.Get(raw, "json_response.things").Array() {
				things = append(things, Thing{
					ThingID:   t.Get("thing_id").String(),
					ThingName: t.Get("thing_metadata.thing_name").String(),
					SecretKey: strings.TrimSuffix(t.Get("thing_metadata.thing_secret_key").String(), "\n"),
				})
			}

			renderTemplate(w, "devices.gohtml", TemplateData{
				Raw:    raw,
				Things: things,
			})
			return
		} else {
			log.Printf("got different challengeName: %v\n", resp.ChallengeName)
			w.WriteHeader(http.StatusInternalServerError)
			renderError(w, fmt.Sprintf("got different challengeName: %v", resp.ChallengeName))
			return
		}
	}
}
