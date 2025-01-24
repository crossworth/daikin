package daikin

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/crossworth/daikin/types"
	"github.com/valyala/fasthttp"
)

// Client creates a Daikin client.
type Client struct {
	target    *url.URL
	secretKey []byte
}

// NewClient creates a *Client.
func NewClient(target *url.URL, secretKey []byte) *Client {
	return &Client{
		target:    target,
		secretKey: secretKey,
	}
}

// makes a request to the given path returning the response as a byte slice.
func (c *Client) makeRequest(ctx context.Context, method string, path string, body []byte) ([]byte, error) {
	// we have to use a custom http client because the server sends invalid http responses
	// example below:
	//  	HTTP/1.1 200 OK
	//  			Content-Type: application/json
	c.target.Path = path
	var (
		endpoint = c.target.String()
		req      = fasthttp.AcquireRequest()
		resp     = fasthttp.AcquireResponse()
	)
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(endpoint)
	req.Header.SetMethod(method)
	req.SetTimeout(30 * time.Second)
	if len(body) > 0 {
		req.SetBody(body)
	}
	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("making request to %q: %w", endpoint, err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request to %q returned status code %d", endpoint, resp.StatusCode())
	}
	data, err := base64.StdEncoding.DecodeString(string(resp.Body()))
	if err != nil {
		return nil, fmt.Errorf("decoding body: %w", err)
	}
	return data, nil
}

// decodeData decodes the given data using the given secretKey.
func decodeData(secretKey []byte, data []byte, path string) ([]byte, error) {
	if len(data) < 16+2 {
		return nil, fmt.Errorf("input data is too short to contain a payload")
	}
	var (
		startPayload = 16
		hasLength    = false
	)
	switch path {
	case "acstatus":
		startPayload = 17
	case "get_scan", "status", "onboard":
		startPayload = 17
		hasLength = true
	}
	var (
		iv         = data[:16]
		length     = data[16] & 255
		payload    = data[startPayload : len(data)-2]
		crc16Bytes = data[len(data)-2:]
		crc16      = (int(crc16Bytes[1]&255) << 8) | int(crc16Bytes[0]&255)
		crc16Check = calculateCRC16(data[:len(data)-2]) & 65535
	)
	_, _ = length, hasLength // not used for now, we don't know what it means or if is required
	if crc16 != crc16Check {
		return nil, fmt.Errorf("invalid crc16")
	}
	decoded, err := decryptAESCFB(payload, secretKey, iv)
	if err != nil {
		return nil, fmt.Errorf("decoding payload, check the secret key")
	}
	return decoded, nil
}

// encodeData encodes the given data using the secretKey.
func encodeData(secretKey []byte, data []byte) ([]byte, error) {
	var (
		iv        = make([]byte, 16)
		inputData = make([]byte, 0, len(data)+2)
	)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("generating IV: %w", err)
	}
	inputData = append(inputData, data...)
	inputData = append(inputData, []byte("BZ")...) // not sure why we should encode the BZ here
	encrypted, err := encryptAESCFB(inputData, secretKey, iv)
	if err != nil {
		return nil, fmt.Errorf("encrypt data: %w", err)
	}
	output := make([]byte, 0, len(encrypted)+18) // 16 bytes iv + 2 bytes crc16
	output = append(output, iv...)
	output = append(output, encrypted...)
	crc16 := calculateCRC16(output)
	output = append(output, byte(crc16&255), byte(crc16>>8)&255)
	return output, nil
}

type StatusResponse struct {
	Username    string `json:"username"`
	StationSSID string `json:"sta_ssid"`
	Status      Status `json:"status"`
}

type Status struct {
	AC    int `json:"ac"`
	STA   int `json:"sta"`
	Cloud int `json:"cloud"`
	Auth  int `json:"auth"`
}

// Status query for the device status.
func (c *Client) Status(ctx context.Context) (*StatusResponse, error) {
	data, err := c.makeRequest(ctx, http.MethodGet, "/status", nil)
	if err != nil {
		return nil, err
	}
	decoded, err := decodeData(c.secretKey, data, "status")
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	var resp StatusResponse
	if err := json.Unmarshal(decoded, &resp); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}
	return &resp, nil
}

type State struct {
	Port1 types.Port `json:"port1"`
	Idu   int        `json:"idu"`
}

// State query for the device state.
func (c *Client) State(ctx context.Context) (*State, error) {
	data, err := c.makeRequest(ctx, http.MethodGet, "/acstatus", nil)
	if err != nil {
		return nil, err
	}
	decoded, err := decodeData(c.secretKey, data, "acstatus")
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	var resp State
	if err := json.Unmarshal(decoded, &resp); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}
	return &resp, nil
}

type DesiredState struct {
	Port1 types.PortState `json:"port1"`
}

// SetState sets the desired state on the device.
func (c *Client) SetState(ctx context.Context, state DesiredState) (*State, error) {
	jsonState, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("encoding json: %w", err)
	}
	stateData, err := encodeData(c.secretKey, jsonState)
	if err != nil {
		return nil, fmt.Errorf("encoding desired state: %w", err)
	}
	data, err := c.makeRequest(ctx, http.MethodPost, "/acstatus", []byte(base64.StdEncoding.EncodeToString(stateData)))
	if err != nil {
		return nil, err
	}
	decoded, err := decodeData(c.secretKey, data, "acstatus")
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	var resp State
	if err := json.Unmarshal(decoded, &resp); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}
	return &resp, nil
}
