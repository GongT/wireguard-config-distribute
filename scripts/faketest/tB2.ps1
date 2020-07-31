#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot/../..

$tmp = [System.IO.Path]::GetTempPath()
New-Item -Path $tmp -Name xxxqqq -ItemType "directory" -Force

$host.ui.RawUI.WindowTitle = "== B2 =="

./dist/client `
	--insecure -D --external-ip-nohttp `
	--hosts-file=$tmp/hosts0 `
	--netgroup=B `
	--server=127.0.0.1 `
	--external-ip=172.0.2.2 `
	--internal-ip=127.0.2.2 `
	--perfer-ip=222.2 `
	--hostname=peer-b2 `
	--title="test B 2" `
	--netgroup="B"
