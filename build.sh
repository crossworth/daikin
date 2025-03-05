#!/usr/bin/env sh

set -eu

env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./build/daikin-windows-32.exe ./cmd/daikin/
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/daikin-windows-64.exe ./cmd/daikin/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./build/daikin-linux-arm ./cmd/daikin/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/daikin-linux-arm64 ./cmd/daikin/
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./build/daikin-linux-386 ./cmd/daikin/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/daikin-linux-amd64 ./cmd/daikin/
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/daikin-darwin-amd64 ./cmd/daikin/


env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./build/ac-server-http-service-windows-32.exe ./cmd/ac-server-http-service/
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/ac-server-http-service-windows-64.exe ./cmd/ac-server-http-service/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./build/ac-server-http-service-linux-arm ./cmd/ac-server-http-service/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/ac-server-http-service-linux-arm64 ./cmd/ac-server-http-service/
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./build/ac-server-http-service-linux-386 ./cmd/ac-server-http-service/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/ac-server-http-service-linux-amd64 ./cmd/ac-server-http-service/
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/ac-server-http-service-darwin-amd64 ./cmd/ac-server-http-service/


env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./build/find-devices-windows-32.exe ./cmd/find-devices/
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/find-devices-windows-64.exe ./cmd/find-devices/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ./build/find-devices-linux-arm ./cmd/find-devices/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/find-devices-linux-arm64 ./cmd/find-devices/
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./build/find-devices-linux-386 ./cmd/find-devices/
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/find-devices-linux-amd64 ./cmd/find-devices/
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/find-devices-darwin-amd64 ./cmd/find-devices/