#!/usr/bin/env pwsh

param([Parameter(Mandatory = $true)][string]$type) 

Set-Location $PSScriptRoot/..

./scripts/build.ps1 $type
if ( $? -eq $false ) { exit 1 }

Copy-Item -v ./dist/$type /usr/local/bin/wireguard-config-$type
