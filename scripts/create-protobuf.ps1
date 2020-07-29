#!/usr/bin/env pwsh

cd $PSScriptRoot/..

New-Item -Name internal/protocol -ItemType "directory" -Force | Out-Null
protoc -I protocol protocol/config-service.proto --go_out=plugins=grpc:internal/protocol
