module github.com/slamy-solutions/openbp/modules/tools/testing

go 1.18

replace github.com/slamy-solutions/openbp/modules/system/libs/golang => ../../system/libs/golang

replace github.com/slamy-solutions/openbp/modules/native/libs/golang => ../../native/libs/golang

replace github.com/slamy-solutions/openbp/modules/tools/libs/sdk/golang => ../libs/sdk/golang

require (
	github.com/stretchr/testify v1.8.2
	google.golang.org/grpc v1.54.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230127162408-596548ed4efa // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
