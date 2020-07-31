#!/usr/bin/env pwsh

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot/../..

./scripts/run.ps1 server --server-name=test --ip-nohttp --debug
