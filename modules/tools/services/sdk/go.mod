module github.com/slamy-solutions/openbp/modules/tools/services/sdk

go 1.18

replace github.com/slamy-solutions/openbp/modules/system/libs/golang => ../../../system/libs/golang

replace github.com/slamy-solutions/openbp/modules/native/libs/golang => ../../../native/libs/golang

replace github.com/slamy-solutions/openbp/modules/tools/libs/sdk/golang => ../../libs/sdk/golang

require (
	github.com/sirupsen/logrus v1.9.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.40.0
	google.golang.org/grpc v1.54.0
)

require (
	cloud.google.com/go/compute v1.18.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230127162408-596548ed4efa // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
