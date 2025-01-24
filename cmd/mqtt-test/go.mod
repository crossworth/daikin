module github.com/crossworth/daikin/cmd/mqtt-test

go 1.23.4

require (
	github.com/crossworth/daikin/mqtt v0.0.0-00010101000000-000000000000
	github.com/eclipse/paho.mqtt.golang v1.5.0 // indirect
)

replace (
	github.com/crossworth/daikin/aws => ../../aws
	github.com/crossworth/daikin/mqtt => ../../mqtt
	github.com/crossworth/daikin/types => ../../types
)

require github.com/crossworth/daikin/aws v0.0.0-00010101000000-000000000000

require (
	github.com/alexrudd/cognito-srp/v4 v4.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.54 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.28.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.49.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.9 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/crossworth/daikin/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)
