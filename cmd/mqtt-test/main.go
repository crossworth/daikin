package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crossworth/daikin/aws"
	"github.com/crossworth/daikin/mqtt"
	"github.com/crossworth/daikin/types"
)

func main() {
	time.Local = time.UTC
	var (
		username string
		password string
		thingID  string
	)
	flag.StringVar(&username, "username", "", "The username")
	flag.StringVar(&password, "password", "", "The password")
	flag.StringVar(&thingID, "thingID", "", "The ThingID")
	flag.Parse()
	if username == "" || password == "" || thingID == "" {
		log.Fatalf("required field not provided\n")
		return
	}

	appContext, appContextCancel := context.WithCancel(context.Background())
	defer appContextCancel()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		appContextCancel()
	}()

	accountInfo, err := aws.GetAccountInfo(appContext, username, password)
	if err != nil {
		log.Printf("failed to get account info: %v\n", err)
		os.Exit(1)
		return
	}

	mqttInfo, err := aws.GetMQTTInfo(appContext, accountInfo.IDToken)
	if err != nil {
		log.Printf("failed to get MQTT info: %v\n", err)
		os.Exit(1)
		return
	}

	client, err := mqtt.NewClient(mqttInfo.AccessKeyID, mqttInfo.SecretKey, mqttInfo.SessionToken)
	if err != nil {
		log.Printf("failed to create client: %v\n", err)
		os.Exit(1)
		return
	}
	defer client.Disconnect()

	printState := func(ctx context.Context, thingID string) {
		thingState, err := client.State(ctx, thingID)
		if err != nil {
			log.Printf("failed to get state: %v\n", err)
			os.Exit(1)
			return
		}
		thingStateBytes, err := json.MarshalIndent(thingState, "", "  ")
		if err != nil {
			log.Printf("failed to marshal state: %v\n", err)
			os.Exit(1)
			return
		}
		fmt.Printf("ThingState:\n\n")
		fmt.Printf("%s\n\n", string(thingStateBytes))
	}

	printState(appContext, thingID)

	fmt.Printf("Set AC to 25°C/Coandă effect/Economy Mode\n")

	if err := client.SetState(appContext, thingID, types.PortState{
		Power:       ptr(1),
		Temperature: ptr(25.0),
		Coanda:      ptr(1),
		Econo:       ptr(1),
	}); err != nil {
		log.Printf("failed to set state: %v\n", err)
		os.Exit(1)
		return
	}

	printState(appContext, thingID)

	fmt.Printf("CTRL+C to stop\n")
	<-appContext.Done()
}

func ptr[T any](v T) *T {
	return &v
}
