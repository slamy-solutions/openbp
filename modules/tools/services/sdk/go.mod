module github.com/slamy-solutions/openbp/modules/tools/services/sdk

go 1.20

replace github.com/slamy-solutions/openbp/modules/system/libs/golang => ../../../system/libs/golang

replace github.com/slamy-solutions/openbp/modules/native/libs/golang => ../../../native/libs/golang

replace github.com/slamy-solutions/openbp/modules/tools/libs/sdk/golang => ../../libs/sdk/golang

require (
	github.com/sirupsen/logrus v1.9.3
	github.com/slamy-solutions/openbp/modules/native/libs/golang v0.0.0-00010101000000-000000000000
	github.com/slamy-solutions/openbp/modules/system/libs/golang v0.0.0-00010101000000-000000000000
	github.com/slamy-solutions/openbp/modules/tools/libs/sdk/golang v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.42.0
	google.golang.org/grpc v1.56.1
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v0.20.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
