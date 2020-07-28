#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

cd $PSScriptRoot/..

$host.ui.RawUI.WindowTitle = "== B1 =="

./dist/client `
--insecure -D `
	--hosts-file=/tmp/xxxqqq/hosts0 `
	--netgroup=B `
	--server=127.0.0.1 `
	--external-ip=127.0.2.1 `
	--internal-ip=127.0.2.1 `
	--perfer-ip=222.1 `
	--hostname=peer-b1 `
	--title="test B 1" `
	--netgroup="B"
