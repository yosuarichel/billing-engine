#!/usr/bin/env bash
RUN_NAME="billing-engine"

mkdir -p output/bin output/conf output/conf/.idl
cp script/* output/
cp conf/* output/conf/
# cp -r conf/.idl/* output/conf/.idl/
chmod +x output/bootstrap.sh

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    echo "generating build for non test env"
    CGO_ENABLED=0 GOOS=linux go build -v -buildvcs=false -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" -o output/bin/${RUN_NAME}
    # go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" -o output/bin/${RUN_NAME}
else
     echo "generating build for test env"
    CGO_ENABLED=0 GOOS=linux go test -c -covermode=set -o output/bin/${RUN_NAME} -coverpkg=./...
    # go test -c -covermode=set -o output/bin/${RUN_NAME} -coverpkg=./...
fi