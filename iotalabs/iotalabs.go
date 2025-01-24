package iotalabs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	iotalabsEndpoint  = "https://dmb.iotalabs.co.in/devices/thinginfo/managething"
	iotalabsUserAgent = "okhttp/5.0.0-alpha.2"
)

// ManageThing returns the output of calling [iotalabsEndpoint].
func ManageThing(ctx context.Context, username string, accessToken string, idToken string) (string, error) {
	body := strings.NewReader(`{"request_type":"GET_THING_INFO","json_request":{},"user_name":"` + username + `","access_token":"` + accessToken + `"}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, iotalabsEndpoint, body)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", iotalabsUserAgent)
	req.Header.Set("Authorization", idToken) // very strange using idToken here and access token on the body
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalidResponse statusCode=%s body=%s", resp.Status, string(responseBody))
	}
	return string(responseBody), nil
}
