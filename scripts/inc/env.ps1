# . "$PSScriptRoot/inc/env.ps1"

Set-Location $PSScriptRoot/../..

Function Test-CommandExists {
	Param ($command)

	try { if (Get-Command $command) { return $true } }

	Catch { return $false }
} 

function IsWindows() {
	$iswin = Get-Variable IsWindows -Scope Global -ErrorAction SilentlyContinue
	if ($iswin ) {
		return $iswin 
	} else {
		return Test-Path "C:\Windows\System32"
	}
}

$T = "github.com/gongt/wireguard-config-distribute/internal/tools"
$GH = git log -1 --pretty=format:%h
git diff-index --quiet HEAD
if ( $? -eq $false ) { $GH += "[M]" }
$env:LDFLAGS += " -X '$T.build_date=$(Get-Date -Format "yyyy-MM-dd+HH:mm:ss")' -X '$T.build_git_hash=$GH'"

if ( ( Test-CommandExists go ) -And ( $(go env GOOS) -Eq "windows" ) ) {
	# go build -ldflags -H=windowsgui -o dist/$type ./cmd/wireguard-config-$type
	# $env:GOGCCFLAGS += " -ldflags -H=windowsgui"
	$env:GOEXE = ".exe"
}
