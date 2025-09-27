#!/usr/bin/env bash
RUN_NAME="billing_engine"

mkdir -p output/bin output/conf output/conf/.idl
cp script/* output/
cp conf/* output/conf/
# cp -r conf/.idl/* output/conf/.idl/
chmod +x output/bootstrap.sh

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    echo "generating build for non test env"
    GOOS=linux GOARCH=arm64 go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" -o output/bin/${RUN_NAME}
else
     echo "generating build for test env"
    GOOS=linux GOARCH=arm64 go test -c -covermode=set -o output/bin/${RUN_NAME} -coverpkg=./...
fi