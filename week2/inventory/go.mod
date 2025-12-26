module github.com/mbakhodurov/homeworks/week2/inventory

go 1.25.4

replace github.com/mbakhodurov/homeworks/week2/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/mbakhodurov/homeworks/week2/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.3.0 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251213004720-97cd9d5aeac2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
)
