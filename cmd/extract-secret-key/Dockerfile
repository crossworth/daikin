FROM golang:1.23-bullseye AS build
ARG SERVICE_VERSION=0
WORKDIR /app
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
COPY . /app
RUN GOOS=linux GOARCH=amd64 go build -o service /app/

FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=build /app/service /app/service
CMD ["/app/service"]