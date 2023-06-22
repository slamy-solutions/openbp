module github.com/slamy-solutions/openbp/modules/tools/testing

go 1.20

replace github.com/slamy-solutions/openbp/modules/system/libs/golang => ../../system/libs/golang

replace github.com/slamy-solutions/openbp/modules/native/libs/golang => ../../native/libs/golang

replace github.com/slamy-solutions/openbp/modules/tools/libs/sdk/golang => ../libs/sdk/golang

require github.com/stretchr/testify v1.8.3

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
