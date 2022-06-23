export PATH="$PATH:$(go env GOPATH)/bin"

# namespace
mkdir -p ./namespace/src/grpc/native_namespace
protoc --go_out=./namespace/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./namespace/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto

# file
mkdir -p ./file/src/grpc/native_namespace
protoc --go_out=./file/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./file/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./file/src/grpc/native_file
protoc --go_out=./file/src/grpc/native_file --go_opt=paths=source_relative --go-grpc_out=./file/src/grpc/native_file --go-grpc_opt=paths=source_relative -I ../proto file.proto