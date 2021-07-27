#!/usr/bin/env pwsh
Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

$Repo = "GongT/wireguard-config-distribute"
$distFolder = "C:/Program Files/WireGuard"
$getFile = "client.exe"
$versionFile = "$distFolder/client.exe.version.txt"
$proxyServer = 'http://proxy-server.:3271/'

. ./lib.ps1

function doUpdate() {
	param($downloadUrl)
	$downloadFile = "$distFolder/$getFile.update"
	$binaryFile = "$distFolder/$getFile"

	Remove-Item $binaryFile -ErrorAction SilentlyContinue | Out-Null
	Remove-Item $downloadFile -ErrorAction SilentlyContinue | Out-Null

	Write-Host -ForegroundColor Gray "    远程: $downloadUrl"
	Write-Host -ForegroundColor Gray "    本地:   $downloadFile"
	Invoke-WebRequest-Wrap -Uri $downloadUrl -OutFile $downloadFile

	stopAllService
	Write-Host "    重命名文件……"
	Copy-Item -Path $downloadFile -Destination $binaryFile
	startAllService
}

detectVersionChange -Repo $Repo -versionFile $versionFile -callback {
	param($change, $downloadUrl)
	if ($change) {
		Write-Host " * 开始下载："
		doUpdate $downloadUrl
	}
	
	Write-Host "当前版本: $(getLocalVersion)"
}
