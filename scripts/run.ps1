#!/usr/bin/env pwsh

param([string]$type)

$ErrorActionPreference = "Stop"
. "$PSScriptRoot/inc/env.ps1"

Clear-Host

$env:GRPC_GO_LOG_SEVERITY_LEVEL = "debug"
$env:GRPC_GO_LOG_VERBOSITY_LEVEL = "99"
$env:GRPC_VERBOSITY = "info"
$env:GRPC_TRACE = "tcp,http,api"

Write-Output "Creating protocol..."
./scripts/create-protobuf.ps1
if ( $? -eq $false ) { exit 1 }

Write-Output "generate..."
go generate ./cmd/wireguard-config-$type
if ( $? -eq $false ) { exit 1 }

Write-Output "run..."
Write-Output ""
Write-Output ""
Write-Output ""
go run ./cmd/wireguard-config-$type @args
