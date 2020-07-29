#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"
cd $PSScriptRoot/../..

$tmp = [System.IO.Path]::GetTempPath()
New-Item -Path $tmp -Name xxxqqq -ItemType "directory" -Force

$host.ui.RawUI.WindowTitle = "== B1 =="

./dist/client `
	--insecure -D --external-ip-nohttp `
	--hosts-file=$tmp/hosts0 `
	--netgroup=B `
	--server=127.0.0.1 `
	--external-ip=172.0.2.1 `
	--internal-ip=127.0.2.1 `
	--perfer-ip=222.1 `
	--hostname=peer-b1 `
	--title="test B 1" `
	--netgroup="B"
