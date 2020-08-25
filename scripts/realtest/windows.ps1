#!/usr/bin/env pwsh
$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot/../..

if ($env:OneDriveConsumer) {
	$Root = "$env:OneDriveConsumer/Software/WireguardConfig"
} elseif ($env:OneDrive) {
	$Root = "$env:OneDrive/Software/WireguardConfig"
} else {
	Write-Error "木有找到 OneDrive 路径"
	Exit-PSSession 1
}

$hashTable = Get-Content -Encoding utf8 "$Root/$env:COMPUTERNAME.conf" | ConvertFrom-StringData
foreach ($key in $hashTable.Keys) {
	$value = $hashTable.$key
	Set-Item env:$key $value 
}
$env:WIREGUARD_LOG = ""

./scripts/run.ps1 client --group normal
