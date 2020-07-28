#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

cd $PSScriptRoot/..

mkdir -p /tmp/xxxqqq/

$host.ui.RawUI.WindowTitle = "== A2 =="

./dist/client `
	--insecure -D --external-ip-nohttp `
	--hosts-file=/tmp/xxxqqq/hosts0 `
	--netgroup=A `
	--server=127.0.0.1 `
	--external-ip=172.0.1.2 `
	--internal-ip=127.0.1.2 `
	--perfer-ip=111.2 `
	--hostname=peer-a2 `
	--title="test A 2" `
	--netgroup="A"
