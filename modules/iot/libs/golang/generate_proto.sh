export PATH="$PATH:$(go env GOPATH)/bin"

# device
echo "Generating proto for core_device service"
mkdir -p ./core/device
protoc --go_out=./core/device --go_opt=paths=source_relative --go-grpc_out=./core/device --go-grpc_opt=paths=source_relative -I ../../proto/core device.proto
# fleet
echo "Generating proto for core_fleet service"
mkdir -p ./core/fleet
protoc --go_out=./core/fleet --go_opt=paths=source_relative --go-grpc_out=./core/fleet --go-grpc_opt=paths=source_relative -I ../../proto/core fleet.proto
# telemetry
echo "Generating proto for core_telemetry service"
mkdir -p ./core/telemetry
protoc --go_out=./core/telemetry --go_opt=paths=source_relative --go-grpc_out=./core/telemetry --go-grpc_opt=paths=source_relative -I ../../proto/core telemetry.proto
# integration_balena
echo "Generating proto for core_integration_balena service"
mkdir -p ./core/integration/balena
protoc --go_out=./core/integration/balena --go_opt=paths=source_relative --go-grpc_out=./core/integration/balena --go-grpc_opt=paths=source_relative -I ../../proto/core/integration balena.proto