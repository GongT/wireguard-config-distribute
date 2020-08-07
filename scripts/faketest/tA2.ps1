#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot/../..

$tmp = [System.IO.Path]::GetTempPath()
New-Item -Path $tmp -Name xxxqqq -ItemType "directory" -Force | out-null

$host.ui.RawUI.WindowTitle = "== A2 =="

$env:WIREGUARD_PASSWORD = Get-Content ~/.wireguard-config-server/password.txt

./dist/client `
	--insecure -D --external-ip-nohttp --external-ip-noupnp --ipv6only --no-upnp-forward --dry `
	--hosts-file=$tmp/hosts0 `
	--netgroup=A `
	--server=127.0.0.1 `
	--external-ip=172.0.1.2 `
	--internal-ip=127.0.1.2 `
	--perfer-ip=111.2 `
	--hostname=peer-a2 `
	--title="test A 2" `
	--netgroup="A" @args
