package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/crossworth/daikin"
	"github.com/crossworth/daikin/types"
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of daikin: --secretKey=<SECRET_KEY> --targetAddress=<TARGET_ADDRESS> <MODE> [state]\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Mode can be one of:\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t- get\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t- set\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "State should be provided only when working with set.\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Example state:\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t'{\"power\":1,\"temperature\":25}'\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Arguments:\n")
		flag.PrintDefaults()
	}
	var (
		inputSecretKey     string
		inputTargetAddress string
	)
	flag.StringVar(&inputSecretKey, "secretKey", "", "The secret key")
	flag.StringVar(&inputTargetAddress, "targetAddress", "", "The target address")
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
		return
	}
	var (
		mode               = strings.ToLower(flag.Arg(0))
		desiredStateString = strings.ToLower(flag.Arg(1))
	)
	if inputSecretKey == "" {
		fmt.Printf("no secretKey provided\n")
		os.Exit(1)
		return
	}
	if inputTargetAddress == "" {
		fmt.Printf("no targetAddress provided\n")
		os.Exit(1)
		return
	}
	target, err := url.Parse(inputTargetAddress)
	if err != nil {
		fmt.Printf("invalid target address: %v\n", err)
		os.Exit(1)
		return
	}
	secretKey, err := base64.StdEncoding.DecodeString(inputSecretKey)
	if err != nil {
		fmt.Printf("could not decode the secret key: %v\n", err)
		os.Exit(1)
		return
	}

	serviceContext, serviceContextCancel := context.WithCancel(context.Background())
	defer serviceContextCancel()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		serviceContextCancel()
	}()
	var (
		client = daikin.NewClient(target, secretKey)
		enc    = json.NewEncoder(os.Stdout)
	)
	enc.SetIndent("", "  ")
	if mode == "get" {
		state, err := client.State(serviceContext)
		if err != nil {
			fmt.Printf("could not get state: %v\n", err)
			os.Exit(1)
			return
		}
		_ = enc.Encode(state)
		return
	}
	if mode == "set" {
		if len(flag.Args()) < 2 {
			flag.Usage()
			os.Exit(1)
			return
		}
		var portState types.PortState
		if err := json.Unmarshal([]byte(desiredStateString), &portState); err != nil {
			fmt.Printf("could not unmarshal desired state %q: %v\n", desiredStateString, err)
			os.Exit(1)
			return
		}
		state, err := client.SetState(serviceContext, daikin.DesiredState{
			Port1: portState,
		})
		if err != nil {
			fmt.Printf("could not set state: %v\n", err)
			os.Exit(1)
			return
		}
		_ = enc.Encode(state)
		return
	}
	fmt.Printf("invalid mode: %s\n", mode)
	os.Exit(1)
	return
}
