export PATH="$PATH:$(go env GOPATH)/bin"

# service
echo "Generating proto for core_service service"
mkdir -p ./core/service
protoc --go_out=./core/service --go_opt=paths=source_relative --go-grpc_out=./core/service --go-grpc_opt=paths=source_relative -I ../../proto/core service.proto
# ticket
echo "Generating proto for core_ticket service"
mkdir -p ./core/ticket
protoc --go_out=./core/ticket --go_opt=paths=source_relative --go-grpc_out=./core/ticket --go-grpc_opt=paths=source_relative -I ../../proto/core ticket.proto