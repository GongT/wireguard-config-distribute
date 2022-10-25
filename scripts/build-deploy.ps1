#!/usr/bin/env pwsh

param([string][Parameter(Mandatory = $true)]$type) 

Set-Location $PSScriptRoot

pwsh build.ps1 $type

Copy-Item -Path ../dist/$type -Destination /usr/local/libexec/wireguard-config/client
