module github.com/crossworth/daikin/mqtt

go 1.23.4

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
	github.com/crossworth/daikin/types v0.0.0-00010101000000-000000000000
)

replace github.com/crossworth/daikin/types => ../types

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
