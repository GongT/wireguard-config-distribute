#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"

mkdir -p /tmp/xxxqqq/

echo "127.0.0.1 local1
1.1.1.1 some-service
" > /tmp/xxxqqq/hosts1

cd $PSScriptRoot/..

$host.ui.RawUI.WindowTitle = "== A1 =="

./dist/client `
	--insecure -D --external-ip-nohttp `
	--hosts-file=/tmp/xxxqqq/hosts1 `
	--netgroup=A `
	--server=127.0.0.1 `
	--external-ip=127.0.1.1 `
	--internal-ip=127.0.1.1 `
	--perfer-ip=111.1 `
	--hostname=peer-a1 `
	--title="test A 1" `
	--netgroup="A"
