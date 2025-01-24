package mqtt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// getSignedURL creates a signed URL for the WebSocket MQTT server.
// We have to implement our own function, the AWS one for some reason don't seem to work.
func getSignedURL(accessKeyID string, secretKey string, sessionToken string) (string, error) {
	if accessKeyID == "" || secretKey == "" {
		return "", fmt.Errorf("credentials cannot be anonymous or empty")
	}
	var (
		tm              = time.Now().UTC()
		amzDate         = tm.UTC().Format("20060102T150405Z")
		dateStamp       = tm.UTC().Format("20060102")
		credentialScope = fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, awsRegion, awsService)
		query           = url.Values{}
	)
	query.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	query.Set("X-Amz-Credential", fmt.Sprintf("%s/%s", accessKeyID, credentialScope))
	query.Set("X-Amz-Date", amzDate)
	query.Set("X-Amz-SignedHeaders", "host")
	var (
		canonicalHeaders       = fmt.Sprintf("host:%s\n", mqttHost)
		hashedPayload          = sha256Hex([]byte(""))
		canonicalRequest       = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", http.MethodGet, "/mqtt", query.Encode(), canonicalHeaders, "host", hashedPayload)
		hashedCanonicalRequest = sha256Hex([]byte(canonicalRequest))
		stringToSign           = fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", amzDate, credentialScope, hashedCanonicalRequest)
		signingKey             = getSigningKey(secretKey, dateStamp, awsRegion, awsService)
		signature              = hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))
	)
	query.Set("X-Amz-Signature", signature)
	if sessionToken != "" {
		query.Set("X-Amz-Security-Token", sessionToken)
	}
	return fmt.Sprintf("wss://%s/mqtt?%s", mqttHost, query.Encode()), nil
}

func sha256Hex(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func getSigningKey(secret, dateStamp, region, service string) []byte {
	data := hmacSHA256([]byte("AWS4"+secret), []byte(dateStamp))
	data = hmacSHA256(data, []byte(region))
	data = hmacSHA256(data, []byte(service))
	return hmacSHA256(data, []byte("aws4_request"))
}
