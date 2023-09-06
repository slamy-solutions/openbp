module github.com/slamy-solutions/openbp/modules/tools/testing

go 1.21.0

replace github.com/slamy-solutions/openbp/modules/system/libs/golang => ../../system/libs/golang

replace github.com/slamy-solutions/openbp/modules/native/libs/golang => ../../native/libs/golang

replace github.com/slamy-solutions/openbp/modules/tools/libs/golang => ../libs/golang

require (
	github.com/slamy-solutions/openbp/modules/tools/libs/golang v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.3
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.56.1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
