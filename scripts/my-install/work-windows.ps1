#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot/../..

if (! (Test-Path $env:GOPATH/bin) ) {
	Write-Error "No GOPATH, or GOPATH/bin did not exists"
	return
}

$hashTable = Get-Content -Encoding utf8 $env:GOPATH/wireguard-client.conf | ConvertFrom-StringData
foreach ($key in $hashTable.Keys) {
	$value = $hashTable.$key
	Set-Item env:$key $value 
}

./scripts/build.ps1 client

$binFile = "$env:GOPATH/bin/wireguard-config-service.exe"

Write-Output ""

if (Test-Path $binFile) {
	Write-Output "Uninstall old service.."
	& $binFile /D /uninstall
} else {
	Write-Output "Old service did not exists."
}
Write-Output "Copy binary file..."
Copy-Item ./dist/client.exe $binFile

Write-Output "Install new service..."
& $binFile /D /install
