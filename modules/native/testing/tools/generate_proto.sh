
function generate {
    # $1 - target (output) folder
    # $2 - source folder
    # $3,4,5,6,7 - source file relative to folder
    echo "Generating proto from $1 to $2"
    protoc \
        --plugin=../../node_modules/.bin/protoc-gen-ts_proto \
        --ts_proto_out=$1 \
        --ts_proto_opt=esModuleInterop=true \
        --ts_proto_opt=env=node \
        -I $2 \
        $3 $4 $5 $6 $7
}

generate ./namespace/proto ../../../modules/native/namespace namespace.proto
