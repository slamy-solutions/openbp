export PATH="$PATH:$(go env GOPATH)/bin"

# namespace
echo "Generating proto for namespace service"
mkdir -p ./namespace
protoc --go_out=./namespace --go_opt=paths=source_relative --go-grpc_out=./namespace --go-grpc_opt=paths=source_relative -I ../../proto namespace.proto
# keyvaluestorage
echo "Generating proto for keyvaluestorage service"
mkdir -p ./keyvaluestorage
protoc --go_out=./keyvaluestorage --go_opt=paths=source_relative --go-grpc_out=./keyvaluestorage --go-grpc_opt=paths=source_relative -I ../../proto keyvaluestorage.proto

# iam_token
echo "Generating proto for iam_token service"
mkdir -p ./iam/token
protoc --go_out=./iam/token --go_opt=paths=source_relative --go-grpc_out=./iam/token --go-grpc_opt=paths=source_relative -I ../../proto/iam token.proto
# iam_policy
echo "Generating proto for iam_policy service"
mkdir -p ./iam/policy
protoc --go_out=./iam/policy --go_opt=paths=source_relative --go-grpc_out=./iam/policy --go-grpc_opt=paths=source_relative -I ../../proto/iam policy.proto
# iam_role
echo "Generating proto for iam_role service"
mkdir -p ./iam/role
protoc --go_out=./iam/role --go_opt=paths=source_relative --go-grpc_out=./iam/role --go-grpc_opt=paths=source_relative -I ../../proto/iam role.proto
# iam_auth
echo "Generating proto for iam_auth service"
mkdir -p ./iam/auth
protoc --go_out=./iam/auth --go_opt=paths=source_relative --go-grpc_out=./iam/auth --go-grpc_opt=paths=source_relative -I ../../proto/iam auth.proto
# iam_config
echo "Generating proto for iam_config service"
mkdir -p ./iam/config
protoc --go_out=./iam/config --go_opt=paths=source_relative --go-grpc_out=./iam/config --go-grpc_opt=paths=source_relative -I ../../proto/iam config.proto
# iam_identity
echo "Generating proto for iam_identity service"
mkdir -p ./iam/identity
protoc --go_out=./iam/identity --go_opt=paths=source_relative --go-grpc_out=./iam/identity --go-grpc_opt=paths=source_relative -I ../../proto/iam identity.proto
# iam_authentication_password
echo "Generating proto for iam_authentication_password service"
mkdir -p ./iam/authentication/password
protoc --go_out=./iam/authentication/password --go_opt=paths=source_relative --go-grpc_out=./iam/authentication/password --go-grpc_opt=paths=source_relative -I ../../proto/iam/authentication password.proto
# iam_authentication_x509
echo "Generating proto for iam_authentication_x509 service"
mkdir -p ./iam/authentication/x509
protoc --go_out=./iam/authentication/x509 --go_opt=paths=source_relative --go-grpc_out=./iam/authentication/x509 --go-grpc_opt=paths=source_relative -I ../../proto/iam/authentication x509.proto
# iam_authentication_oauth2
echo "Generating proto for iam_authentication_oauth2 service"
mkdir -p ./iam/authentication/oauth2
protoc --go_out=./iam/authentication/oauth2 --go_opt=paths=source_relative --go-grpc_out=./iam/authentication/oauth2 --go-grpc_opt=paths=source_relative -I ../../proto/iam/authentication oauth2.proto
# iam_actor_user
echo "Generating proto for iam_actor_user service"
mkdir -p ./iam/actor/user
protoc --go_out=./iam/actor/user --go_opt=paths=source_relative --go-grpc_out=./iam/actor/user --go-grpc_opt=paths=source_relative -I ../../proto/iam/actor user.proto

# storage_bucket
echo "Generating proto for storage_bucket service"
mkdir -p ./storage/bucket
protoc --go_out=./storage/bucket --go_opt=paths=source_relative --go-grpc_out=./storage/bucket --go-grpc_opt=paths=source_relative -I ../../proto/storage bucket.proto
# storage_fs
echo "Generating proto for storage_fs service"
mkdir -p ./storage/fs
protoc --go_out=./storage/fs --go_opt=paths=source_relative --go-grpc_out=./storage/fs --go-grpc_opt=paths=source_relative -I ../../proto/storage fs.proto