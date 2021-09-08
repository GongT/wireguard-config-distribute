#!/usr/bin/env pwsh

param([string]$type) 

# $ErrorActionPreference = "Stop"
. "$PSScriptRoot/inc/env.ps1"
. "$PSScriptRoot/inc/x.ps1"

$iargs = @()

Set-Location $PSScriptRoot/..
New-Item -ItemType Directory -Force dist | Out-Null

function build() {
	param (
		[parameter(position = 0, Mandatory = $true)][string]$type,
		[parameter(Mandatory = $false)][switch]$docker
	)
	
	Write-Output "Generate $type..."
	& x go generate ./cmd/wireguard-config-$type

	Write-Output "Build $type${env:GOEXE}..."

	$out = "$type${env:GOEXE}"
	if ($env:CI) {
		Write-Host "::set-output name=artifact::$out"
	}

	$verb = @()
	if ($env:CI) {
		# $verb = @('-x', '-v')
	}
	
	[string[]]$build = @( '-ldflags', $env:LDFLAGS) + $iargs + $args + @('-o', "dist/$out", "./cmd/wireguard-config-$type")
	x 'go' 'build' @verb @build
}

if ($env:RUNNER_TEMP) {
	$env:TMP = $env:RUNNER_TEMP
	$env:TEMP = $env:RUNNER_TEMP
} else {
	$TMP = New-Item -Name ".temp" -ItemType Directory -Force
	$env:TMP = $TMP.FullName
	$env:TEMP = $TMP.FullName
}

if ($type -Eq "musl") {
	SetExecuteMethod -container
	x go env
} elseif ($env:CI) {
	Write-Output "=============================================="
	Get-ChildItem Env:* | Format-Table
	Write-Output "=============================================="
	go env
	Write-Output "=============================================="
} else {
	Clear-Host
}

# Set-PSDebug -Trace 1

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
	} elseif ($type -Eq "musl") {
		# & $env:CC -v
		$env:GOOS = "linux"
		$env:GOARCH = "arm64"
		$env:GOEXE = ".alpine"

		build server
		build client
		build tool
	} else {
		build $type
	}
} elseif ($IsWindows) {
	build client
	build tool
} else {
	build server
	build client
	build tool
}
