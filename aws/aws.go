package aws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &userAgentTransport{
		UserAgent: awsUserAgent,
		Base:      http.DefaultTransport,
	},
}

type AccountInfoResponse struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
	IDToken     string `json:"idToken"`
}

// GetAccountInfo gets the account info for the given username and password.
func GetAccountInfo(ctx context.Context, username string, password string) (*AccountInfoResponse, error) {
	// configure cognito srp
	csrp, err := cognitosrp.NewCognitoSRP(username, password, poolID, clientID, aws.String(clientSecret))
	if err != nil {
		return nil, fmt.Errorf("creating SRP: %v", err)
	}

	// configure cognito identity provider
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithHTTPClient(httpClient),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %v", err)
	}

	svc := cip.NewFromConfig(cfg)

	// initialize the auth
	resp, err := svc.InitiateAuth(ctx, &cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		return nil, fmt.Errorf("calling InitiateAuth: %v", err)
	}

	// we only support PASSWORD_VERIFIER challenge
	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, err := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())
		if err != nil {
			return nil, fmt.Errorf("calling PasswordVerifierChallenge: %v", err)
		}

		resp, err := svc.RespondToAuthChallenge(ctx, &cip.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			return nil, fmt.Errorf("calling RespondToAuthChallenge: %v", err)
		}

		if resp.AuthenticationResult.AccessToken == nil || resp.AuthenticationResult.IdToken == nil {
			return nil, fmt.Errorf("accessToken or idToken is nil")
		}

		// get the user
		user, err := svc.GetUser(ctx, &cip.GetUserInput{
			AccessToken: resp.AuthenticationResult.AccessToken,
		})
		if err != nil {
			return nil, fmt.Errorf("calling GetUser: %v", err)
		}

		// got nil username
		if user.Username == nil {
			return nil, fmt.Errorf("user.Username is nil")
		}

		return &AccountInfoResponse{
			Username:    *user.Username,
			AccessToken: *resp.AuthenticationResult.AccessToken,
			IDToken:     *resp.AuthenticationResult.IdToken,
		}, nil
	}
	return nil, fmt.Errorf("got different challengeName: %v", resp.ChallengeName)
}

type MQTTInfoResponse struct {
	AccessKeyID  string `json:"accessKeyID"`
	SecretKey    string `json:"secretKey"`
	SessionToken string `json:"sessionToken"`
}

// GetMQTTInfo gets the mqtt info for the given idToken and user.
func GetMQTTInfo(ctx context.Context, idToken string) (*MQTTInfoResponse, error) {
	// configure cognito identity provider
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithHTTPClient(httpClient),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %v", err)
	}

	// create the client
	client := cognitoidentity.NewFromConfig(cfg)

	logins := map[string]string{
		"cognito-idp." + awsRegion + ".amazonaws.com/" + poolID: idToken,
	}

	getIDResponse, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String(identityPoolID),
		Logins:         logins,
	})
	if err != nil {
		return nil, fmt.Errorf("calling GetId: %v", err)
	}

	credentialsResponse, err := client.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: getIDResponse.IdentityId,
		Logins:     logins,
	})
	if err != nil {
		return nil, fmt.Errorf("calling GetCredentialsForIdentity: %v", err)
	}

	return &MQTTInfoResponse{
		AccessKeyID:  *credentialsResponse.Credentials.AccessKeyId,
		SecretKey:    *credentialsResponse.Credentials.SecretKey,
		SessionToken: *credentialsResponse.Credentials.SessionToken,
	}, nil
}
