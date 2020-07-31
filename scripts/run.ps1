#!/usr/bin/env pwsh

param([string]$type)

Set-Location $PSScriptRoot/..

Clear-Host

$env:GRPC_GO_LOG_SEVERITY_LEVEL = "debug"
$env:GRPC_GO_LOG_VERBOSITY_LEVEL = "99"
$env:GRPC_VERBOSITY = "info"
$env:GRPC_TRACE = "tcp,http,api"

Write-Output "Creating protocol..."
./scripts/create-protobuf.ps1

Write-Output "generate..."
go generate ./cmd/$type/*.go

Write-Output "run..."
Write-Output ""
Write-Output ""
Write-Output ""
go run ./cmd/$type @args
