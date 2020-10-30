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
	Write-Host -ForegroundColor Gray "    远程: $downloadUrl"
	Write-Host -ForegroundColor Gray "    本地:   $downloadFile"
	Invoke-WebRequest-Wrap -Uri $downloadUrl -Out $downloadFile

	$serviceConfigList = Get-ChildItem -Filter *.xml
	foreach ($item in $serviceConfigList ) {
		Write-Host "    停止服务：$item"
		./winsw.exe stop $item
	}
	Write-Host "    重命名文件……"
	Copy-Item -Path $downloadFile -Destination $binaryFile
	foreach ($item in $serviceConfigList ) {
		Write-Host "    启动服务：$item"
		./winsw.exe start $item
	}
}

detectVersionChange -Repo $Repo -versionFile $versionFile -callback {
	param($change, $downloadUrl)
	if ($change) {
		Write-Host " * 有更新，开始下载："
		doUpdate $downloadUrl
	} else {
		Write-Host " * 已是最新版本"
		Write-Host -ForegroundColor Gray "    文件:   $binaryFile"
		return $binaryFile
	}
}
