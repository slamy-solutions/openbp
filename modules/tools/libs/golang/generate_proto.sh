export PATH="$PATH:$(go env GOPATH)/bin"

### NATIVE ###
rm -rf ./native

# namespace
echo "Generating proto for namespace service"
mkdir -p ./native/namespace
protoc --go_out=./native/namespace --go_opt=paths=source_relative --go-grpc_out=./native/namespace --go-grpc_opt=paths=source_relative -I ../../proto/sdk/native namespace.proto

### SDK ###
rm -rf ./sdk

echo "Generating proto for sdk service"
mkdir -p ./sdk/sdk
protoc --go_out=./sdk/sdk --go_opt=paths=source_relative --go-grpc_out=./sdk/sdk --go-grpc_opt=paths=source_relative -I ../../proto/sdk/sdk sdk.proto