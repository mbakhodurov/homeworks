module github.com/mbakhodurov/homeworks/week5/notification

go 1.25.4

replace github.com/mbakhodurov/homeworks/week5/shared => ../shared

replace github.com/mbakhodurov/homeworks/week5/platform => ../platform

require (
	github.com/IBM/sarama v1.47.0
	github.com/caarlos0/env/v11 v11.4.0
	github.com/go-telegram/bot v1.20.0
	github.com/mbakhodurov/homeworks/week5/platform v0.0.0-00010101000000-000000000000
	github.com/mbakhodurov/homeworks/week5/shared v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.25 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
