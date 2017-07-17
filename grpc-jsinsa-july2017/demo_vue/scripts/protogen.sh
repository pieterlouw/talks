#!/bin/bash


protoc --plugin=protoc-gen-ts=../node_modules/.bin/protoc-gen-ts -I ../proto --js_out=import_style=commonjs,binary:../client/_proto  --ts_out=service=true:../client/_proto ../proto/carservice.proto