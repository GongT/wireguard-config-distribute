#!/usr/bin/env pwsh

param([string]$type) 

# $ErrorActionPreference = "Stop"
. "$PSScriptRoot/inc/env.ps1"
. "$PSScriptRoot/inc/x.ps1"

$iargs = @()

function build() {
	param (
		[parameter(position = 0, Mandatory = $true)][string]$type,
		[parameter(Mandatory = $false)][switch]$docker
	)
	
	Write-Output "Generate $type..."
	& x go generate ./cmd/wireguard-config-$type
	if ( $? -eq $false ) { exit 1 }

	Write-Output "Build $type${env:GOEXE}..."

	$out = "$type${env:GOEXE}"
	Write-Host "::set-output name=artifact::$out"

	[string[]]$build = @('go', 'build', '-x', '-v', '-ldflags', $env:LDFLAGS) + $iargs + $args + @('-o', "dist/$out", "./cmd/wireguard-config-$type")
	& x @build
	if ( $? -eq $false ) {
		Write-Output "Failed go build..."
		Start-Sleep -Seconds 5
		Write-Output "Quit with 1..."
		exit 1 
	}
}

if ($env:CI) {
	if ($env:RUNNER_TEMP) {
		$env:TMP = $env:RUNNER_TEMP
		$env:TEMP = $env:RUNNER_TEMP
	} else {
		$TMP = New-Item -Name ".temp" -ItemType Directory -Force
		$env:TMP = $TMP.FullName
		$env:TEMP = $TMP.FullName
	}

	Write-Output "=============================================="
	Get-ChildItem Env:*
	Write-Output "=============================================="
	go env
	Write-Output "=============================================="
} else {
	Clear-Host
}

# Set-PSDebug -Trace 1

if ($type -Eq "musl") {
	# & $env:CC -v
	$env:GOOS = "linux"
	$env:GOARCH = "arm64"
	$env:GOEXE = ".alpine"
	$type = "client"

	SetExecuteMethod -container
}

Write-Output "Creating protocol..."
. ./scripts/create-protobuf.ps1
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
	} else {
		build $type
	}
} else {
	build server
	build client
	build tool
}
