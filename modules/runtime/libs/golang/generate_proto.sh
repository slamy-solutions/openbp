export PATH="$PATH:$(go env GOPATH)/bin"

# manager_runtime
echo "Generating proto for manager_runtime service"
mkdir -p ./manager/runtime
protoc --go_out=./manager/runtime --go_opt=paths=source_relative --go-grpc_out=./manager/runtime --go-grpc_opt=paths=source_relative -I ../../proto/manager runtime.proto

# manager_rpc
echo "Generating proto for manager_rpc service"
mkdir -p ./manager/rpc
protoc --go_out=./manager/rpc --go_opt=paths=source_relative --go-grpc_out=./manager/rpc --go-grpc_opt=paths=source_relative -I ../../proto/manager rpc.proto