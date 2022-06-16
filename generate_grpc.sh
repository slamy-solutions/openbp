export PATH="$PATH:$(go env GOPATH)/bin"

# catalog
protoc --go_out=./catalog/grpc/native_catalog_grpc --go_opt=paths=source_relative --go-grpc_out=./catalog/grpc/native_catalog_grpc --go-grpc_opt=paths=source_relative -I ./proto proto/catalog.proto
protoc --go_out=./catalog/grpc/native_cache_grpc --go_opt=paths=source_relative --go-grpc_out=./catalog/grpc/native_cache_grpc --go-grpc_opt=paths=source_relative -I ./proto proto/cache.proto

#cache
protoc --go_out=./cache/grpc/native_cache_grpc --go_opt=paths=source_relative --go-grpc_out=./cache/grpc/native_cache_grpc --go-grpc_opt=paths=source_relative -I ./proto proto/cache.proto

#lambda
protoc --plugin=/usr/local/bin/protoc-gen-ts_proto --ts_proto_out=./lambda/src/proto --ts_proto_opt=esModuleInterop=true -I ./proto ./proto/catalog.proto