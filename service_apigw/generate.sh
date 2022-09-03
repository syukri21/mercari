#!/bin/bash

THEGOPATH=$(go env GOPATH)
GOMODCACHE=$(go env GOMODCACHE)
GRPC_GATEWAY_VERSION=2.10.0

for i in "$@"
do :
    echo "generating ${i} service"
    
    echo "  - grpc bindings"
    protoc -I/usr/local/include -I. \
    -I$THEGOPATH/src \
    -I$GOMODCACHE/github.com/grpc-ecosystem/grpc-gateway/v2@v$GRPC_GATEWAY_VERSION \
    --go_out=. \
    --go-grpc_out=. \
    protos/${i}.proto
    
    echo "  - grpc-gateway"
    protoc -I/usr/local/include -I. \
    -I$THEGOPATH/src \
    -I$GOMODCACHE/github.com/grpc-ecosystem/grpc-gateway/v2@v$GRPC_GATEWAY_VERSION \
    --grpc-gateway_out=logtostderr=true,grpc_api_configuration=protos/${i}.yml:. \
    protos/${i}.proto
    
    # move generated sources
    mkdir -p internal/generated/${i}
    mv protos/${i}*.go internal/generated/${i}/.
done

echo "gRPC stubs generated"
