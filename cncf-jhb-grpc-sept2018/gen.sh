#!/bin/bash

set -e

cd $(dirname $0)/..

rm -rf generated/*

# generate the protobufs
protoc --go_out=plugins=grpc:./gen -I./proto ./proto/kv.proto