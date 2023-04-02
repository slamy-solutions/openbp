export PATH="$PATH:$(go env GOPATH)/bin"

# vault
echo "Generating proto for vault service"
mkdir -p ./vault
protoc --go_out=./vault --go_opt=paths=source_relative --go-grpc_out=./vault --go-grpc_opt=paths=source_relative -I ../../proto vault.proto
