#!/usr/bin/env pwsh

param([string]$type)

cd $PSScriptRoot/..

Clear-Host

$env:GRPC_GO_LOG_SEVERITY_LEVEL = "debug"
$env:GRPC_GO_LOG_VERBOSITY_LEVEL = "99"
$env:GRPC_VERBOSITY = "info"
$env:GRPC_TRACE = "tcp,http,api"

echo "Creating protocol..."
./scripts/create-protobuf.ps1

echo "generate..."
go generate ./cmd/$type/*.go

echo "run..."
echo ""
echo ""
echo ""
go run ./cmd/$type @args
