#!/usr/bin/env pwsh

cd $PSScriptRoot/..

./scripts/create-protobuf.ps1

echo "building server..."
go build -o dist/server ./cmd/server

echo "building tool..."
go build -o dist/tool ./cmd/tool

echo "building client..."
go build -o dist/client ./cmd/client
