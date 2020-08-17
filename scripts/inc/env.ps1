# . "$PSScriptRoot/inc/env.ps1"
Set-Location $PSScriptRoot/../..

$T = "github.com/gongt/wireguard-config-distribute/internal/tools"
$GH = git log -1 --pretty=format:%h
git diff-index --quiet HEAD
if ( $? -eq $false ) { $GH += " (has modify)" }
$env:LDFLAGS += " -X '$T.build_date=$(Get-Date -Format "yyyy/MM/dd+HH:mm:ss")' -X '$T.build_git_hash=$GH'"

$ext = $env:GOEXE
if ( $(go env GOOS) -Eq "windows" ) {
	# go build -ldflags -H=windowsgui -o dist/$type ./cmd/wireguard-config-$type
	# $env:GOGCCFLAGS += " -ldflags -H=windowsgui"
	$ext = ".exe"
}
