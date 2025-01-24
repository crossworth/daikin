package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crossworth/daikin/types"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

const (
	mqttHost   = "a2zl28lcjjddi3-ats.iot.us-east-1.amazonaws.com"
	awsRegion  = "us-east-1"
	awsService = "iotdata"
)

// Client is a MTQClient for the device.
type Client struct {
	client mqtt.Client
}

// NewClient returns a new client connection to the server.
func NewClient(accessKeyID string, secretKey string, sessionToken string) (*Client, error) {
	broker, err := getSignedURL(accessKeyID, secretKey, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("creating signed URL: %w", err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(uuid.New().String())
	opts.SetHTTPHeaders(http.Header{
		"Host": []string{},
	})
	opts.SetCleanSession(true)
	var (
		client = mqtt.NewClient(opts)
		token  = client.Connect()
	)
	token.Wait()
	if token.Error() != nil {
		return nil, fmt.Errorf("connecting to server: %w", token.Error())
	}
	return &Client{
		client: client,
	}, nil
}

// IsConnected returns if the client is connected.
func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}

// Disconnect disconnects from the server.
func (c *Client) Disconnect() {
	c.client.Disconnect(0)
}

// pubSub is a helper function used to publish to a topic and get a response.
func (c *Client) pubSub(ctx context.Context, publishTopic string, receiveTopic string, data string) ([]byte, error) {
	// we add 10 seconds of timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	responseChan := make(chan []byte, 1)
	token := c.client.Subscribe(receiveTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		select {
		case responseChan <- msg.Payload():
			client.Unsubscribe(receiveTopic)
		default:
		}
	})
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("subscribing to topic: %w", err)
	}
	token = c.client.Publish(publishTopic, 0, false, data)
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("publishing to topic: %w", err)
	}
	select {
	case res := <-responseChan:
		return res, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

type Message struct {
	State State `json:"state"`
}

type State struct {
	Reported *ReportedState `json:"reported,omitempty"`
	Desired  *DesiredState  `json:"desired,omitempty"`
}
type ReportedState struct {
	Connected  int        `json:"connected"`
	Src        int        `json:"src"`
	Port1      types.Port `json:"port1"`
	ACUnitType string     `json:"ac_unit_type"`
}

type DesiredState struct {
	Port1 types.PortState `json:"port1"`
	Src   int             `json:"src"`
}

// State gets the current device state.
func (c *Client) State(ctx context.Context, thingID string) (*ReportedState, error) {
	var (
		publishTopic = "$aws/things/" + thingID + "/shadow/get"
		receiveTopic = "$aws/things/" + thingID + "/shadow/get/accepted"
	)
	resp, err := c.pubSub(ctx, publishTopic, receiveTopic, "")
	if err != nil {
		return nil, err
	}
	var message Message
	if err := json.Unmarshal(resp, &message); err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}
	return message.State.Reported, nil
}

// SetState sets the current device state, if a field in [PortState] is nil, we will not change the value.
func (c *Client) SetState(ctx context.Context, thingID string, state types.PortState) error {
	var (
		publishTopic = "$aws/things/" + thingID + "/shadow/update"
		receiveTopic = "$aws/things/" + thingID + "/shadow/update/accepted"
	)
	desiredState, err := json.Marshal(Message{
		State: State{
			Desired: &DesiredState{
				Port1: state,
				Src:   5, // not sure why 5, but the app always sends a hardcoded 5
			},
		},
	})
	if err != nil {
		return fmt.Errorf("marshalling state: %w", err)
	}
	if _, err := c.pubSub(ctx, publishTopic, receiveTopic, string(desiredState)); err != nil {
		return err
	}
	return nil
}
