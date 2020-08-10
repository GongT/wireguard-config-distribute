#!/usr/bin/env pwsh

param([string]$type) 

Set-Location $PSScriptRoot/..

if (($env:SystemRoot) -And (Test-Path "$env:SystemRoot")) {
	# go build -ldflags -H=windowsgui -o dist/$type ./cmd/$type
	# $env:GOGCCFLAGS += " -ldflags -H=windowsgui"
	$ext = ".exe"
}

function build() {
	param ([Parameter(Mandatory)]$type)
	
	Write-Output "Generate $type..."
	go generate ./cmd/$type
	if ( $? -eq $false ) { exit 1 }

	Write-Output "Build $type..."
	go build @args -o dist/$type$ext ./cmd/$type
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
