package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	time.Local = time.UTC

	serviceContext, serviceContextCancel := context.WithCancel(context.Background())
	defer serviceContextCancel()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		serviceContextCancel()
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", indexHandler)
	mux.HandleFunc("POST /info", accountInfoHandler())

	var (
		httpServer = &http.Server{
			ReadHeaderTimeout: 20 * time.Second,
			ReadTimeout:       1 * time.Minute,
			WriteTimeout:      2 * time.Minute,
			IdleTimeout:       120 * time.Second,
			Addr:              ":8080",
			Handler:           mux,
			ErrorLog:          log.Default(),
		}
		httpServerErrors = make(chan error, 1)
	)

	log.Printf("stating http server at http://127.0.0.1:8080\n")
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
