#!/usr/bin/env pwsh

param([string]$type) 

$GOPATH = $(go env GOPATH)

$ErrorActionPreference = "Stop"
. "$PSScriptRoot/inc/env.ps1"

$iargs = @()

function x() {
	param (
		[string]$cmd,
		[Parameter(ValueFromRemainingArguments)][string[]]$args
	)
	$argstr = $args -join " "
	Write-Host -Separator " " -ForegroundColor Gray $cmd $argstr
	& $cmd @args
	if ( $? -eq $false ) { exit 1 }
}
function build() {
	param (
		[parameter(position = 0, Mandatory = $true)][string]$type,
		[parameter(Mandatory = $false)][switch]$docker
	)
	
	Write-Output "Generate $type..."
	go generate ./cmd/wireguard-config-$type
	if ( $? -eq $false ) { exit 1 }

	Write-Output "Build $type${env:GOEXE}..."

	[string[]]$build = @('go', 'build', '-ldflags', $env:LDFLAGS) + $iargs + $args + @('-o', "dist/$type${env:GOEXE}", "./cmd/wireguard-config-$type")
	if ($docker) {
		x podman run --rm `
			"--workdir=/data" `
			"--volume=$(Get-Location):/data" `
			"--volume=${GOPATH}:/go" `
			"golang:alpine" `
			@build
	} else {
		x @build
	}
}

Clear-Host

Write-Output "Creating protocol..."
./scripts/create-protobuf.ps1
if ( $? -eq $false ) { exit 1 }

if ($type) {
	if ($type -Eq "android") {
		# & $env:CC -v
		$env:GOOS = "linux"
		$env:GOARCH = "arm64"
		$env:GOEXE = ".android"
		$env:GOARM = 7
		$env:CGO_ENABLED = 0
		$iargs += "-tags", "moveable"
		build client
	} elseif ($type -Eq "musl") {
		# & $env:CC -v
		$env:GOOS = "linux"
		$env:GOARCH = "arm64"
		$env:GOEXE = ".alpine"
		build client -docker
	} else {
		build $type
	}
} else {
	build server
	build client
	build tool
}
