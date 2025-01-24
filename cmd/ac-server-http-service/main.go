package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crossworth/daikin"
	"github.com/crossworth/daikin/types"
)

func main() {
	var (
		inputSecretKey     string
		inputTargetAddress string
		listenAddress      string
	)
	flag.StringVar(&inputSecretKey, "secretKey", "", "The secret key")
	flag.StringVar(&inputTargetAddress, "targetAddress", "", "The target address")
	flag.StringVar(&listenAddress, "listenAddress", ":8080", "The listen addresss")
	flag.Parse()
	if inputSecretKey == "" {
		log.Fatalf("no secretKey provided\n")
		return
	}
	if inputTargetAddress == "" {
		log.Fatalf("no targetAddress provided\n")
		return
	}
	target, err := url.Parse(inputTargetAddress)
	if err != nil {
		log.Fatalf("invalid target address: %v\n", err)
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
	secretKey, err := base64.StdEncoding.DecodeString(inputSecretKey)
	if err != nil {
		log.Fatalf("could not decode the secret key: %v\n", err)
		return
	}

	client := daikin.NewClient(target, secretKey)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		writer.Header().Set("Content-Type", "application/json")

		state, err := client.State(request.Context())
		if err != nil {
			log.Printf("could not get AC status: %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}{
				Error:   false,
				Message: "could not get AC status: " + err.Error(),
			})
			return
		}
		writer.WriteHeader(http.StatusOK)
		_ = encoder.Encode(state)
		return
	})
	mux.HandleFunc("POST /state", func(writer http.ResponseWriter, request *http.Request) {
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		writer.Header().Set("Content-Type", "application/json")

		var portState types.PortState
		if err := json.NewDecoder(request.Body).Decode(&portState); err != nil {
			log.Printf("could not decode the request body: %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}{
				Error:   false,
				Message: "could not decode the request body: " + err.Error(),
			})
			return
		}
		state, err := client.SetState(serviceContext, daikin.DesiredState{
			Port1: portState,
		})
		if err != nil {
			log.Printf("could not set state: %v\n", err)
			writer.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}{
				Error:   false,
				Message: "could not set state: " + err.Error(),
			})
			return
		}
		writer.WriteHeader(http.StatusOK)
		_ = encoder.Encode(state)
		return
	})

	var (
		httpServer = &http.Server{
			ReadHeaderTimeout: 20 * time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      2 * time.Minute,
			IdleTimeout:       120 * time.Second,
			Addr:              listenAddress,
			Handler:           mux,
			ErrorLog:          log.Default(),
		}
		httpServerErrors = make(chan error, 1)
	)

	log.Printf("starting http server at %s\n", listenAddress)
	go func() {
		httpServerErrors <- httpServer.ListenAndServe()
	}()

	select {
	case <-serviceContext.Done():
	case err := <-httpServerErrors:
		log.Printf("http server error: %v\n", err)
	}
	log.Printf("stopping service\n")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown error: %v\n", err)
	}
}
