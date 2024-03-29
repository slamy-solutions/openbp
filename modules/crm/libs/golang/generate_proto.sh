export PATH="$PATH:$(go env GOPATH)/bin"

# client
echo "Generating proto for core_client service"
mkdir -p ./core/client
protoc --go_out=./core/client --go_opt=paths=source_relative --go-grpc_out=./core/client --go-grpc_opt=paths=source_relative -I ../../proto client.proto
# settings
echo "Generating proto for settings service"
mkdir -p ./core/settings
protoc --go_out=./core/settings --go_opt=paths=source_relative --go-grpc_out=./core/settings --go-grpc_opt=paths=source_relative -I ../../proto settings.proto
# onecsync
echo "Generating proto for onecsync service"
mkdir -p ./core/onecsync
protoc --go_out=./core/onecsync --go_opt=paths=source_relative --go-grpc_out=./core/onecsync --go-grpc_opt=paths=source_relative -I ../../proto onecsync.proto
# performer
echo "Generating proto for performer service"
mkdir -p ./core/performer
protoc --go_out=./core/performer --go_opt=paths=source_relative --go-grpc_out=./core/performer --go-grpc_opt=paths=source_relative -I ../../proto performer.proto
# department
echo "Generating proto for department service"
mkdir -p ./core/department
protoc --go_out=./core/department --go_opt=paths=source_relative --go-grpc_out=./core/department --go-grpc_opt=paths=source_relative -I ../../proto department.proto
# project
echo "Generating proto for project service"
mkdir -p ./core/project
protoc --go_out=./core/project --go_opt=paths=source_relative --go-grpc_out=./core/project --go-grpc_opt=paths=source_relative -I ../../proto project.proto
# kanban
echo "Generating proto for kanban service"
mkdir -p ./core/kanban
protoc --go_out=./core/kanban --experimental_allow_proto3_optional --go_opt=paths=source_relative --go-grpc_out=./core/kanban --go-grpc_opt=paths=source_relative -I ../../proto kanban.proto