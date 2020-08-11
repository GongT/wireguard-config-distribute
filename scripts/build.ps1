#!/usr/bin/env pwsh

param([string]$type) 

Set-Location $PSScriptRoot/..

if ( $(go env GOOS) -Eq "windows" ) {
	# go build -ldflags -H=windowsgui -o dist/$type ./cmd/wireguard-config-$type
	# $env:GOGCCFLAGS += " -ldflags -H=windowsgui"
	$ext = ".exe"
}

function build() {
	param ([Parameter(Mandatory)]$type)
	
	Write-Output "Generate $type..."
	go generate ./cmd/wireguard-config-$type
	if ( $? -eq $false ) { exit 1 }

	Write-Output "Build $type..."
	go build @args -o dist/$type$ext ./cmd/wireguard-config-$type
	if ( $? -eq $false ) { exit 1 }
}

clear-host

Write-Output "Creating protocol..."
./scripts/create-protobuf.ps1
if ( $? -eq $false ) { exit 1 }

if ($type) {
	build $type
} else {
	build server
	build client
	build tool
}
