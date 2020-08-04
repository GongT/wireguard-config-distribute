#!/usr/bin/env pwsh

Set-Location $PSScriptRoot/../..


$hashTable = Get-Content -Encoding utf8 $env:GOPATH/wireguard-client.conf | ConvertFrom-StringData
foreach ($key in $hashTable.Keys) {
	$value = $hashTable.$key
	Set-Item env:$key $value 
}

./scripts/run.ps1 client /install
