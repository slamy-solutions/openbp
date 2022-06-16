export PATH="$PATH:$(go env GOPATH)/bin"

# namespace
mkdir -p grpc/native_namespace
protoc --go_out=./grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./grpc/native_namespace --go-grpc_opt=paths=source_relative namespace.proto