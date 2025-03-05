package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crossworth/daikin"
)

func main() {
	var (
		timeout      time.Duration
		outputInJSON bool
	)
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Timeout")
	flag.BoolVar(&outputInJSON, "json", false, "Output in JSON format")
	flag.Parse()
	if timeout <= 0 {
		timeout = 30
	}

	serviceContext, serviceContextCancel := context.WithCancel(context.Background())
	defer serviceContextCancel()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		serviceContextCancel()
	}()

	devices, err := daikin.DiscoveryDevices(serviceContext, timeout)
	if err != nil {
		fmt.Printf("error discovering devices: %v\n", err)
		os.Exit(1)
		return
	}

	if outputInJSON {
		data, err := json.MarshalIndent(devices, "", "  ")
		if err != nil {
			fmt.Printf("error encoding json: %v\n", err)
			os.Exit(1)
			return
		}
		fmt.Println(string(data))
	} else {
		if len(devices) == 0 {
			fmt.Printf("no device found\n")
			return
		}
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}
}
