export PATH="$PATH:$(go env GOPATH)/bin"

# namespace
mkdir -p ./namespace/src/grpc/native_namespace
protoc --go_out=./namespace/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./namespace/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto

# file
mkdir -p ./file/src/grpc/native_namespace
protoc --go_out=./file/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./file/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./file/src/grpc/native_file
protoc --go_out=./file/src/grpc/native_file --go_opt=paths=source_relative --go-grpc_out=./file/src/grpc/native_file --go-grpc_opt=paths=source_relative -I ../proto file.proto

# lambda/manager
mkdir -p ./lambda/manager/src/grpc/native_namespace
protoc --go_out=./lambda/manager/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./lambda/manager/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./lambda/manager/src/grpc/native_lambda
protoc --go_out=./lambda/manager/src/grpc/native_lambda --go_opt=paths=source_relative --go-grpc_out=./lambda/manager/src/grpc/native_lambda --go-grpc_opt=paths=source_relative -I ../proto lambda.proto

#lambda-entrypoint
mkdir -p ./lambda/entrypoint/src/grpc/native_namespace
protoc --go_out=./lambda/entrypoint/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./lambda/entrypoint/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./lambda/entrypoint/src/grpc/native_lambda
protoc --go_out=./lambda/entrypoint/src/grpc/native_lambda --go_opt=paths=source_relative --go-grpc_out=./lambda/entrypoint/src/grpc/native_lambda --go-grpc_opt=paths=source_relative -I ../proto lambda.proto

#keyvaluestorage
mkdir -p ./keyvaluestorage/src/grpc/native_namespace
protoc --go_out=./keyvaluestorage/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./keyvaluestorage/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./keyvaluestorage/src/grpc/native_keyvaluestorage
protoc --go_out=./keyvaluestorage/src/grpc/native_keyvaluestorage --go_opt=paths=source_relative --go-grpc_out=./keyvaluestorage/src/grpc/native_keyvaluestorage --go-grpc_opt=paths=source_relative -I ../proto keyvaluestorage.proto

#iam-config
mkdir -p ./iam/config/src/grpc/native_iam_configuration
protoc --go_out=./iam/config/src/grpc/native_iam_configuration --go_opt=paths=source_relative --go-grpc_out=./iam/config/src/grpc/native_iam_configuration --go-grpc_opt=paths=source_relative -I ../proto/iam configuration.proto
mkdir -p ./iam/config/src/grpc/native_keyvaluestorage
protoc --go_out=./iam/config/src/grpc/native_keyvaluestorage --go_opt=paths=source_relative --go-grpc_out=./iam/config/src/grpc/native_keyvaluestorage --go-grpc_opt=paths=source_relative -I ../proto keyvaluestorage.proto

#iam-policy
mkdir -p ./iam/policy/src/grpc/native_namespace
protoc --go_out=./iam/policy/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./iam/policy/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./iam/policy/src/grpc/native_iam_policy
protoc --go_out=./iam/policy/src/grpc/native_iam_policy --go_opt=paths=source_relative --go-grpc_out=./iam/policy/src/grpc/native_iam_policy --go-grpc_opt=paths=source_relative -I ../proto/iam policy.proto