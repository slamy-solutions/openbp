export PATH="$PATH:$(go env GOPATH)/bin"

# namespace
echo "Generating proto for namespace service"
mkdir -p ./namespace/src/grpc/native_namespace
protoc --go_out=./namespace/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./namespace/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto

# file
echo "Generating proto for file service"
mkdir -p ./file/src/grpc/native_namespace
protoc --go_out=./file/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./file/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./file/src/grpc/native_file
protoc --go_out=./file/src/grpc/native_file --go_opt=paths=source_relative --go-grpc_out=./file/src/grpc/native_file --go-grpc_opt=paths=source_relative -I ../proto file.proto

# lambda/manager
echo "Generating proto for lambda_manager service"
mkdir -p ./lambda/manager/src/grpc/native_namespace
protoc --go_out=./lambda/manager/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./lambda/manager/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./lambda/manager/src/grpc/native_lambda
protoc --go_out=./lambda/manager/src/grpc/native_lambda --go_opt=paths=source_relative --go-grpc_out=./lambda/manager/src/grpc/native_lambda --go-grpc_opt=paths=source_relative -I ../proto lambda.proto

# lambda-entrypoint
echo "Generating proto for lambda_entrypoint service"
mkdir -p ./lambda/entrypoint/src/grpc/native_namespace
protoc --go_out=./lambda/entrypoint/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./lambda/entrypoint/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./lambda/entrypoint/src/grpc/native_lambda
protoc --go_out=./lambda/entrypoint/src/grpc/native_lambda --go_opt=paths=source_relative --go-grpc_out=./lambda/entrypoint/src/grpc/native_lambda --go-grpc_opt=paths=source_relative -I ../proto lambda.proto

# keyvaluestorage
echo "Generating proto for keyvaluestorage service"
mkdir -p ./keyvaluestorage/src/grpc/native_namespace
protoc --go_out=./keyvaluestorage/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./keyvaluestorage/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./keyvaluestorage/src/grpc/native_keyvaluestorage
protoc --go_out=./keyvaluestorage/src/grpc/native_keyvaluestorage --go_opt=paths=source_relative --go-grpc_out=./keyvaluestorage/src/grpc/native_keyvaluestorage --go-grpc_opt=paths=source_relative -I ../proto keyvaluestorage.proto

# iam-config
echo "Generating proto for iam_config service"
mkdir -p ./iam/config/src/grpc/native_iam_configuration
protoc --go_out=./iam/config/src/grpc/native_iam_configuration --go_opt=paths=source_relative --go-grpc_out=./iam/config/src/grpc/native_iam_configuration --go-grpc_opt=paths=source_relative -I ../proto/iam configuration.proto
mkdir -p ./iam/config/src/grpc/native_keyvaluestorage
protoc --go_out=./iam/config/src/grpc/native_keyvaluestorage --go_opt=paths=source_relative --go-grpc_out=./iam/config/src/grpc/native_keyvaluestorage --go-grpc_opt=paths=source_relative -I ../proto keyvaluestorage.proto

# iam-policy
echo "Generating proto for iam_policy service"
mkdir -p ./iam/policy/src/grpc/native_namespace
protoc --go_out=./iam/policy/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./iam/policy/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./iam/policy/src/grpc/native_iam_policy
protoc --go_out=./iam/policy/src/grpc/native_iam_policy --go_opt=paths=source_relative --go-grpc_out=./iam/policy/src/grpc/native_iam_policy --go-grpc_opt=paths=source_relative -I ../proto/iam policy.proto

# iam-identity
echo "Generating proto for iam_identity service"
mkdir -p ./iam/identity/src/grpc/native_namespace
protoc --go_out=./iam/identity/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./iam/identity/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./iam/identity/src/grpc/native_iam_policy
protoc --go_out=./iam/identity/src/grpc/native_iam_policy --go_opt=paths=source_relative --go-grpc_out=./iam/identity/src/grpc/native_iam_policy --go-grpc_opt=paths=source_relative -I ../proto/iam policy.proto
mkdir -p ./iam/identity/src/grpc/native_iam_identity
protoc --go_out=./iam/identity/src/grpc/native_iam_identity --go_opt=paths=source_relative --go-grpc_out=./iam/identity/src/grpc/native_iam_identity --go-grpc_opt=paths=source_relative -I ../proto/iam identity.proto

# iam-token
echo "Generating proto for iam_token service"
mkdir -p ./iam/token/src/grpc/native_namespace
protoc --go_out=./iam/token/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./iam/token/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./iam/token/src/grpc/native_iam_token
protoc --go_out=./iam/token/src/grpc/native_iam_token --go_opt=paths=source_relative --go-grpc_out=./iam/token/src/grpc/native_iam_token --go-grpc_opt=paths=source_relative -I ../proto/iam token.proto

# iam-authentication-password
echo "Generating proto for iam_authentication_password service"
mkdir -p ./iam/authentication/password/src/grpc/native_namespace
protoc --go_out=./iam/authentication/password/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./iam/authentication/password/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./iam/authentication/password/src/grpc/native_iam_authentication_password
protoc --go_out=./iam/authentication/password/src/grpc/native_iam_authentication_password --go_opt=paths=source_relative --go-grpc_out=./iam/authentication/password/src/grpc/native_iam_authentication_password --go-grpc_opt=paths=source_relative -I ../proto/iam/authentication password.proto

# iam-oauth
echo "Generating proto for iam_oauth service"
mkdir -p ./iam/oauth/src/grpc/native_namespace
protoc --go_out=./iam/oauth/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./iam/oauth/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./iam/oauth/src/grpc/native_iam_policy
protoc --go_out=./iam/oauth/src/grpc/native_iam_policy --go_opt=paths=source_relative --go-grpc_out=./iam/oauth/src/grpc/native_iam_policy --go-grpc_opt=paths=source_relative -I ../proto/iam policy.proto
mkdir -p ./iam/oauth/src/grpc/native_iam_identity
protoc --go_out=./iam/oauth/src/grpc/native_iam_identity --go_opt=paths=source_relative --go-grpc_out=./iam/oauth/src/grpc/native_iam_identity --go-grpc_opt=paths=source_relative -I ../proto/iam identity.proto
mkdir -p ./iam/oauth/src/grpc/native_iam_oauth
protoc --go_out=./iam/oauth/src/grpc/native_iam_oauth --go_opt=paths=source_relative --go-grpc_out=./iam/oauth/src/grpc/native_iam_oauth --go-grpc_opt=paths=source_relative -I ../proto/iam oauth.proto
mkdir -p ./iam/oauth/src/grpc/native_iam_token
protoc --go_out=./iam/oauth/src/grpc/native_iam_token --go_opt=paths=source_relative --go-grpc_out=./iam/oauth/src/grpc/native_iam_token --go-grpc_opt=paths=source_relative -I ../proto/iam token.proto
mkdir -p ./iam/oauth/src/grpc/native_iam_authentication_password
protoc --go_out=./iam/oauth/src/grpc/native_iam_authentication_password --go_opt=paths=source_relative --go-grpc_out=./iam/oauth/src/grpc/native_iam_authentication_password --go-grpc_opt=paths=source_relative -I ../proto/iam/authentication password.proto

# actor-user
echo "Generating proto for actor_user service"
mkdir -p ./actor/user/src/grpc/native_iam_identity
protoc --go_out=./actor/user/src/grpc/native_iam_identity --go_opt=paths=source_relative --go-grpc_out=./actor/user/src/grpc/native_iam_identity --go-grpc_opt=paths=source_relative -I ../proto/iam identity.proto
mkdir -p ./actor/user/src/grpc/native_actor_user
protoc --go_out=./actor/user/src/grpc/native_actor_user --go_opt=paths=source_relative --go-grpc_out=./actor/user/src/grpc/native_actor_user --go-grpc_opt=paths=source_relative -I ../proto/actor user.proto

# api
echo "Generating proto for api service"
mkdir -p ./api/src/grpc/native_namespace
protoc --go_out=./api/src/grpc/native_namespace --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_namespace --go-grpc_opt=paths=source_relative -I ../proto namespace.proto
mkdir -p ./api/src/grpc/native_iam_identity
protoc --go_out=./api/src/grpc/native_iam_identity --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_iam_identity --go-grpc_opt=paths=source_relative -I ../proto/iam identity.proto
mkdir -p ./api/src/grpc/native_iam_policy
protoc --go_out=./api/src/grpc/native_iam_policy --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_iam_policy --go-grpc_opt=paths=source_relative -I ../proto/iam policy.proto
mkdir -p ./api/src/grpc/native_iam_token
protoc --go_out=./api/src/grpc/native_iam_token --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_iam_token --go-grpc_opt=paths=source_relative -I ../proto/iam token.proto
mkdir -p ./api/src/grpc/native_iam_authentication_password
protoc --go_out=./api/src/grpc/native_iam_authentication_password --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_iam_authentication_password --go-grpc_opt=paths=source_relative -I ../proto/iam/authentication password.proto
mkdir -p ./api/src/grpc/native_iam_auth
protoc --go_out=./api/src/grpc/native_iam_auth --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_iam_auth --go-grpc_opt=paths=source_relative -I ../proto/iam oauth.proto
mkdir -p ./api/src/grpc/native_actor_user
protoc --go_out=./api/src/grpc/native_actor_user --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_actor_user --go-grpc_opt=paths=source_relative -I ../proto/actor user.proto
mkdir -p ./api/src/grpc/native_keyvaluestorage
protoc --go_out=./api/src/grpc/native_keyvaluestorage --go_opt=paths=source_relative --go-grpc_out=./api/src/grpc/native_keyvaluestorage --go-grpc_opt=paths=source_relative -I ../proto keyvaluestorage.proto