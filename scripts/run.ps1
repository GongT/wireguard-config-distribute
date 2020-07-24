#!/usr/bin/env pwsh

param([string]$type)

cd $PSScriptRoot/..

Clear-Host

echo "Creating protocol..."
./scripts/create-protobuf.ps1

echo "generate..."
go generate ./cmd/$type/*.go

echo "run..."
go run ./cmd/$type @args
