#!/usr/bin/env pwsh

param([string]$type="tool") 

cd $PSScriptRoot/..

./scripts/build.ps1 -type $type

echo "Copy binary..."
scp dist/$TYPE home.gongt.me:/tmp

echo "Execute..."
ssh -tt home.gongt.me /tmp/$TYPE @args
