#!/usr/bin/env pwsh

Set-Location $PSScriptRoot/../..

$env:WIREGUARD_SERVER = "grpc.services.gongt.me:443"
$env:WIREGUARD_NETWORK = "home"
$env:WIREGUARD_TITLE = "桌面"
$env:WIREGUARD_PUBLIC_IP_NO_HTTP = "true"
$env:WIREGUARD_CONFIG_DEVELOPMENT = "true"
$env:WIREGUARD_REQUEST_IP = "0.50"
$env:WIREGUARD_LOG = "A:/wireguard-output.log"

./scripts/run.ps1 client
