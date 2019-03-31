#!/bin/bash
GOPATH=/c/projects/grpc
PATH=$GOPATH/bin:$PATH
TOOLCHAIN_DIR=$PWD/toolchain
TOOLCHAIN_INCLUDE=$TOOLCHAIN_DIR/include
TOOLCHAIN_BIN=$TOOLCHAIN_DIR/win64/bin
PROJECT_ROOT=$GOPATH
$TOOLCHAIN_BIN/protoc $PWD/proto/echo.proto -I $PWD -I $TOOLCHAIN_INCLUDE/google/protobuf/ --go_out=plugins=grpc:. #$GOPATH/src/
