#!/usr/bin/env pwsh

Set-Location $PSScriptRoot/..
$ErrorActionPreference = "Stop"
. "$PSScriptRoot/inc/x.ps1"

New-Item -Name internal/protocol -ItemType "directory" -Force | Out-Null
x protoc -I protocol protocol/config-service.proto --go_out=plugins=grpc:internal/protocol
