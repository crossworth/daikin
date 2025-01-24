FROM golang:1.23-bullseye AS build
WORKDIR /app
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
COPY . /app
RUN CGO_ENABLED=0 go build -o daikin /app/cmd/daikin
RUN CGO_ENABLED=0 go build -o ac-server-http-service /app/cmd/ac-server-http-service

FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=build /app/daikin /app/daikin
COPY --from=build /app/ac-server-http-service /app/ac-server-http-service