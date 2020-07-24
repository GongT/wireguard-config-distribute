#!/usr/bin/env pwsh

param([string]$type) 

cd $PSScriptRoot/..

function build() {
	param ([Parameter(Mandatory)]$type)
	
	echo "Building $type..."
	go generate ./cmd/$type
	go build -o dist/$type ./cmd/$type
}


echo "Creating protocol..."
./scripts/create-protobuf.ps1

if ($type) {
	build $type
} else{
	build server
	build client
	build tool
}
