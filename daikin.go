package daikin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

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
func (c *Client) makeRequest(ctx context.Context, method string, path string) ([]byte, error) {
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

// decodeResponse decodes the response.
func (c *Client) decodeResponse(data []byte) (json.RawMessage, error) {
	if len(data) < 16+2 {
		return nil, fmt.Errorf("input data is too short to contain a payload")
	}
	var (
		iv         = data[:16]
		payload    = data[17 : len(data)-2]
		crc16Bytes = data[len(data)-2:]
		crc16      = (int(crc16Bytes[1]&255) << 8) | int(crc16Bytes[0]&255)
		crc16Check = calculateCRC16(data[:len(data)-2]) & 65535
	)
	if crc16 != crc16Check {
		return nil, fmt.Errorf("invalid crc16")
	}
	decoded, err := decryptAESCFB(payload, c.secretKey, iv)
	if err != nil {
		return nil, fmt.Errorf("decoding payload, check the secret key")
	}
	return decoded, nil
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
	data, err := c.makeRequest(ctx, http.MethodGet, "/status")
	if err != nil {
		return nil, err
	}
	decoded, err := c.decodeResponse(data)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	var resp StatusResponse
	if err := json.Unmarshal(decoded, &resp); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}
	return &resp, nil
}

type ACStatus struct {
	Port1 Port `json:"port1"`
	Idu   int  `json:"idu"`
}

type Mode int

func (m Mode) String() string {
	switch m {
	case 0:
		return "Automático"
	case 2:
		return "Desumidificar"
	case 3:
		return "Resfriar"
	case 4:
		return "Aquecer"
	case 6:
		return "Ventilar"
	default:
		return fmt.Sprintf("%d", m)
	}
}

type Fan int

func (f Fan) String() string {
	switch f {
	case 3:
		return "Baixa"
	case 4:
		return "Média-Baixa"
	case 5:
		return "Média"
	case 6:
		return "Média-Alta"
	case 7:
		return "Alta"
	case 17:
		return "Automático"
	case 18:
		return "Silencioso"
	default:
		return fmt.Sprintf("%d", f)
	}
}

type Port struct {
	Power         int     `json:"power"`
	Mode          Mode    `json:"mode"`
	Temperature   float64 `json:"temperature"`
	Fan           Fan     `json:"fan"`
	HSwing        int     `json:"h_swing"`
	VSwing        int     `json:"v_swing"`
	Coanda        int     `json:"coanda"`
	Econo         int     `json:"econo"`
	Powerchill    int     `json:"powerchill"`
	GoodSleep     int     `json:"good_sleep"`
	Streamer      int     `json:"streamer"`
	OutQuite      int     `json:"out_quite"`
	OnTimerSet    int     `json:"on_timer_set"`
	OnTimerValue  int     `json:"on_timer_value"`
	OffTimerSet   int     `json:"off_timer_set"`
	OffTimerValue int     `json:"off_timer_value"`
	Sensors       Sensors `json:"sensors"`
	RstR          int     `json:"rst_r"`
	FWVer         string  `json:"fw_ver"`
}

type Sensors struct {
	RoomTemp float64 `json:"room_temp"`
	OutTemp  float64 `json:"out_temp"`
}

// ACStatus query for the ac status.
func (c *Client) ACStatus(ctx context.Context) (*ACStatus, error) {
	data, err := c.makeRequest(ctx, http.MethodGet, "/acstatus")
	if err != nil {
		return nil, err
	}
	decoded, err := c.decodeResponse(data)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	var resp ACStatus
	if err := json.Unmarshal(decoded, &resp); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}
	return &resp, nil
}
