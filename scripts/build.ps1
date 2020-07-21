#!/usr/bin/env pwsh

cd $PSScriptRoot/..

./scripts/create-protobuf.ps1
go build -o dist/server ./cmd/server
go build -o dist/client ./cmd/client
