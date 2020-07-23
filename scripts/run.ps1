#!/usr/bin/env pwsh

param([string]$type) 

cd $PSScriptRoot/..

if (!(Test-Path("internal/protocol"))){
	echo "Creating protocol..."
	./scripts/create-protobuf.ps1
}

echo "generate..."
go generate ./cmd/$type/*.go

echo "run..."
go run cmd/$type/*.go @args
