#!/usr/bin/env pwsh

Set-Location $PSScriptRoot/..
$ErrorActionPreference = "Stop"
. "$PSScriptRoot/inc/x.ps1"

New-Item -Name internal/protocol -ItemType "directory" -Force | Out-Null
x protoc -I protocol -I/usr/include protocol/config-service.proto `
	--go_out=internal/protocol --go_opt=paths=source_relative `
	--go-grpc_out=internal/protocol --go-grpc_opt=paths=source_relative 
