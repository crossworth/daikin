package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crossworth/daikin"
	"github.com/gosuri/uilive"
)

func main() {
	var (
		inputSecretKey     string
		inputTargetAddress string
	)
	flag.StringVar(&inputSecretKey, "secretKey", "", "The secret key")
	flag.StringVar(&inputTargetAddress, "targetAddress", "", "The target address")
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
	var (
		client = daikin.NewClient(target, secretKey)
		writer = uilive.New()
	)
	writer.Start()
	defer writer.Stop()
loop:
	for {
		select {
		case <-serviceContext.Done():
			break loop
		case <-time.After(1 * time.Second):
			acStatusResp, err := client.State(serviceContext)
			if err != nil {
				_, _ = fmt.Fprintf(writer, "Erro ao ler o status do AC: %v\n", err)
			} else {
				_, _ = fmt.Fprintf(writer, "AC Status\n")
				_, _ = fmt.Fprintf(writer.Newline(), "Ligado: %t\n", acStatusResp.Port1.Power == 1)
				_, _ = fmt.Fprintf(writer.Newline(), "Modo: %s\n", acStatusResp.Port1.Mode.String())
				_, _ = fmt.Fprintf(writer.Newline(), "Temperatura: %.2f°C\n", acStatusResp.Port1.Temperature)
				_, _ = fmt.Fprintf(writer.Newline(), "Velocidade: %s\n", acStatusResp.Port1.Fan.String())
				_, _ = fmt.Fprintf(writer.Newline(), "Oscilação: %t\n", acStatusResp.Port1.VSwing == 1)
				_, _ = fmt.Fprintf(writer.Newline(), "Conforto: %t\n", acStatusResp.Port1.Coanda == 1)
				_, _ = fmt.Fprintf(writer.Newline(), "Econômico: %t\n", acStatusResp.Port1.Econo == 1)
				_, _ = fmt.Fprintf(writer.Newline(), "Potente: %t\n", acStatusResp.Port1.Powerchill == 1)
				if acStatusResp.Port1.OnTimerSet == 1 {
					_, _ = fmt.Fprintf(writer.Newline(), "On timer: %d minutos\n", acStatusResp.Port1.OnTimerValue)
				}
				if acStatusResp.Port1.OffTimerSet == 1 {
					_, _ = fmt.Fprintf(writer.Newline(), "Off timer: %d minutos\n", acStatusResp.Port1.OffTimerValue)
				}
				_, _ = fmt.Fprintf(writer.Newline(), "Temperatura ambiente: %.2f°C\n", acStatusResp.Port1.Sensors.RoomTemp)
				_, _ = fmt.Fprintf(writer.Newline(), "Temperatura externa: %.2f°C\n", acStatusResp.Port1.Sensors.OutTemp)
				_, _ = fmt.Fprintf(writer.Newline(), "Versão do firmware: %s\n", acStatusResp.Port1.FWVer)
			}
		}
	}
}
