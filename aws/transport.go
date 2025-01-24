package aws

import (
	"net/http"
)

type userAgentTransport struct {
	UserAgent string
	Base      http.RoundTripper
}

// RoundTrip executes the HTTP request, injecting the custom User-Agent.
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newReq := req.Clone(req.Context())
	newReq.Header.Set("User-Agent", t.UserAgent)
	if t.Base == nil {
		t.Base = http.DefaultTransport
	}
	return t.Base.RoundTrip(newReq)
}
