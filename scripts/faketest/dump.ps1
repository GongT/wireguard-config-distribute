#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot/../..

$env:WIREGUARD_PASSWORD = Get-Content ~/.wireguard-config-server/password.txt

./scripts/run.ps1 client --insecure -s 127.0.0.1 dump
