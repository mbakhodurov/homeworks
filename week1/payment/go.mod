module github.com/mbakhodurov/homeworks/week1/payment

replace github.com/mbakhodurov/homeworks/week1/shared => ../shared

go 1.25.4

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/mbakhodurov/homeworks/week1/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.78.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251111163417-95abcf5c77ba // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
