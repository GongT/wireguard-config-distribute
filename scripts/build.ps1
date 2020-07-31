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
	Write-Output "Build $type..."
	go build @args -o dist/$type$ext ./cmd/$type
}


Write-Output "Creating protocol..."
./scripts/create-protobuf.ps1

if ($type) {
	build $type
} else {
	build server
	build client
	build tool
}
