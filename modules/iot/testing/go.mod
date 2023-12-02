module github.com/slamy-solutions/openbp/modules/iot/testing

go 1.21.4

replace github.com/slamy-solutions/openbp/modules/iot/libs/golang => ../libs/golang

replace github.com/slamy-solutions/openbp/modules/native/libs/golang => ../../native/libs/golang

replace github.com/slamy-solutions/openbp/modules/system/libs/golang => ../../system/libs/golang

require (
	github.com/slamy-solutions/openbp/modules/iot/libs/golang v0.0.0-00010101000000-000000000000
	github.com/slamy-solutions/openbp/modules/native/libs/golang v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.59.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
