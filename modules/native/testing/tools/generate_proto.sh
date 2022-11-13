
function generate {
    # $1 - target (output) folder
    # $2 - source folder
    # $3,4,5,6,7 - source file relative to folder
    echo "Generating proto to $1 from $2"
    mkdir -p $1
    protoc \
        --plugin=../../../../node_modules/.bin/protoc-gen-ts_proto \
        --ts_proto_out=$1 \
        --ts_proto_opt=esModuleInterop=true \
        --ts_proto_opt=env=node \
        -I $2 \
        $3 $4 $5 $6 $7
}

# namespace
generate ./namespace/proto ../../proto namespace.proto
generate ./file/proto ../../proto namespace.proto file.proto
generate ./keyvaluestorage/proto ../../proto keyvaluestorage.proto
generate ./iam/proto ../../proto/iam oauth.proto configuration.proto identity.proto policy.proto token.proto oauth.proto
generate ./iam/proto/authentication ../../proto/iam/authentication password.proto
generate ./actor/proto ../../proto/actor user.proto