#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot/../..

$tmp = [System.IO.Path]::GetTempPath()
New-Item -Path $tmp -Name xxxqqq -ItemType "directory" -Force | out-null

Write-Output "127.0.0.1 local1
1.1.1.1 some-service
" > $tmp/hosts1

$host.ui.RawUI.WindowTitle = "== A1 =="

$env:WIREGUARD_PASSWORD = Get-Content ~/.wireguard-config-server/password.txt

./scripts/build.ps1 client
if ( $? -eq $false ) { exit 1 }

./dist/client `
	--insecure -D --disable-listen --dry `
	--hosts-file=$tmp/hosts1 `
	--netgroup=A `
	--server=127.0.0.1 `
	--perfer-ip=111.1 `
	--hostname=peer-a1 `
	--title="test A 1" `
	--interface="wgt_a1" `
	--machine-id="machineidA1" `
	@args
