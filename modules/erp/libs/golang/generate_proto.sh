export PATH="$PATH:$(go env GOPATH)/bin"

# catalog
echo "Generating proto for core_catalog service"
mkdir -p ./core/catalog
protoc --go_out=./core/catalog --go_opt=paths=source_relative --go-grpc_out=./core/catalog --go-grpc_opt=paths=source_relative -I ../../proto/core catalog.proto
