#!/usr/bin/env pwsh

cd $PSScriptRoot/..

mkdir -p internal/protocol
protoc -I protocol protocol/config-service.proto --go_out=plugins=grpc:internal/protocol
