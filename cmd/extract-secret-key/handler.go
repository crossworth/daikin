package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/crossworth/daikin/aws"
	"github.com/crossworth/daikin/iotalabs"
	"github.com/tidwall/gjson"
)

//go:embed templates/*
var templatesFS embed.FS

type TemplateData struct {
	Error    string                `json:"error"`
	Raw      string                `json:"raw"`
	Things   []Thing               `json:"things"`
	MQTTInfo *aws.MQTTInfoResponse `json:"mqttInfo"`
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

func accountInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	accountInfo, err := aws.GetAccountInfo(ctx, username, password)
	if err != nil {
		log.Printf("error getting account info: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		renderError(w, fmt.Sprintf("error getting account info: %v", err))
		return
	}

	mqttInfo, err := aws.GetMQTTInfo(ctx, accountInfo.IDToken)
	if err != nil {
		log.Printf("error getting mqtt info: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		renderError(w, fmt.Sprintf("error getting mqtt info: %v", err))
		return
	}

	// call the managething endpoint
	raw, err := iotalabs.ManageThing(ctx, accountInfo.Username, accountInfo.AccessToken, accountInfo.IDToken)
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
		Raw:      raw,
		Things:   things,
		MQTTInfo: mqttInfo,
	})
	return
}
