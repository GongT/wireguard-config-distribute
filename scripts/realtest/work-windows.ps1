#!/usr/bin/env pwsh

Set-Location $PSScriptRoot/../..

$env:WIREGUARD_SERVER = "grpc.services.gongt.me:443"
$env:WIREGUARD_NETWORK = "work"
$env:WIREGUARD_TITLE = "工作机(windows)"
$env:WIREGUARD_IPV6 = "true"
$env:WIREGUARD_PUBLIC_IP_NO_UPNP = "true"
$env:WIREGUARD_PUBLIC_IP_NO_HTTP = "true"
$env:WIREGUARD_NO_UPNP = "true"
$env:WIREGUARD_CONFIG_DEVELOPMENT = "true"
$env:WIREGUARD_REQUEST_IP = "1.1"
$env:WIREGUARD_LOG = "D:/Projects/Go/GOPATH/output.log"

go run ./cmd/client /install
