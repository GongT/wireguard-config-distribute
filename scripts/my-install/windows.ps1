#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot/../..

if (! (Test-Path $env:GOPATH/bin) ) {
	Write-Error "No GOPATH, or GOPATH/bin did not exists"
	return
}

./scripts/build.ps1 client
if ( $? -eq $false ) { exit 1 }

if ($env:OneDriveConsumer) {
	$Root = "$env:OneDriveConsumer/Software/WireguardConfig"
} elseif ($env:OneDrive) {
	$Root = "$env:OneDrive/Software/WireguardConfig"
} else {
	Write-Error "木有找到 OneDrive 路径"
	Exit-PSSession 1
}

$binFile = "$Root/wireguard-config-service.exe"

$hashTable = Get-Content -Encoding utf8 "$Root/$env:COMPUTERNAME.conf" | ConvertFrom-StringData
foreach ($key in $hashTable.Keys) {
	$value = $hashTable.$key
	Set-Item env:$key $value 
}

Write-Output ""

function uninstall() {
	Write-Output "Uninstall old service.."
	& $binFile /uninstall
	if ( $? -eq $false ) { exit 1 }
	Start-Sleep -Seconds 5
}

if (Test-Path $binFile) {
	uninstall
	Copy-Item -Force ./dist/client.exe $binFile
} else {
	Copy-Item -Force ./dist/client.exe $binFile
	uninstall
}

Write-Output "Copy binary file..."
Copy-Item -Force ./dist/client.exe $binFile

Write-Output "Install new service..."
& $binFile /install
