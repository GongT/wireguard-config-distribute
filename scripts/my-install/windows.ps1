#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
	Write-Error "必须使用管理员权限运行脚本"
}

Set-Location $PSScriptRoot/../..

function buildGithubReleaseUrl() {
	param (
		[Parameter(Mandatory)][string]$Repo,
		[Parameter(Mandatory)][string]$TagName,
		[Parameter(Mandatory)][string]$GetFile
	)
	return "https://github.com/$Repo/releases/download/$TagName/$GetFile"
}

function downloadGithubRelease() {
	param (
		[Parameter(Mandatory)][string]$Repo,
		[Parameter(Mandatory)][string]$GetFile,
		[Parameter(Mandatory)][string]$SaveAs,
		[Parameter(Mandatory)][string]$DistFolder
	)

	$releaseDataUrl = "https://api.github.com/repos/$Repo/releases?page=1&per_page=1"
	$versionFile = "$DistFolder/$SaveAs.version.txt"
	$binaryFile = "$DistFolder/$SaveAs"
	
	Write-Host "检查 $Repo 版本……"
	Write-Host -ForegroundColor Gray "    来源： $releaseDataUrl"
	$releaseData = (Invoke-WebRequest-Wrap -Uri $releaseDataUrl | ConvertFrom-Json)[0]
	
	if (Test-Path -Path $versionFile) {
		[int]$versionLocal = Get-Content -Encoding utf8 $versionFile
	} else {
		[int]$versionLocal = 0
	}
	if ($versionLocal -eq $releaseData.id ) {
		Write-Host " * 已是最新版本"
		Write-Host -ForegroundColor Gray "    文件:   $binaryFile"
		return $binaryFile
	}

	$downloadUrl = buildGithubReleaseUrl -Repo $Repo -TagName $releaseData.tag_name -GetFile $GetFile
	# $releaseData.assets
	Write-Host " * 有更新，开始下载："
	Write-Host -ForegroundColor Gray "    远程: $downloadUrl"
	Write-Host -ForegroundColor Gray "    本地:   $binaryFile"
	Invoke-WebRequest-Wrap -Uri $downloadUrl -Out $binaryFile 

	Set-Content -Encoding utf8 -Path $versionFile -Value $releaseData.id

	return $binaryFile
}


function createConfig() {
	Write-Output @"
<service>
	<id>wg_$($env:WIREGUARD_GROUP)</id>
	<name>Wireguard 自动配置[组：$($env:WIREGUARD_GROUP)]</name>
	<executable>$clientBinary</executable>
	<download
		failOnError="true"
		from="$downloadUrl"
		to="$clientBinary"
		proxy="$proxyServer"
	/>
	<delayedAutoStart>true</delayedAutoStart>
	<stoptimeout>20sec</stoptimeout>
	<onfailure
		action="restart"
		delay="20sec" />
	<serviceaccount>
		<username>NT AUTHORITY\NetworkService</username>
	</serviceaccount>
	<log mode="roll">
		<logpath>$(Split-Path -Parent $env:WIREGUARD_LOG)</logpath>
	</log>
	<description>自动配置Wireguard网络

配置内容：
"@

	foreach ($key in $hashTable.Keys) {
		if ($key -eq "WIREGUARD_LOG") {
			continue
		}
		$value = $hashTable.$key
		Write-Output "[$key=$value]"
	}
	Write-Output "</description>"
	foreach ($key in $hashTable.Keys) {
		if ($key -eq "WIREGUARD_LOG") {
			continue
		}
		$value = $hashTable.$key
		Write-Output @"
	<env
		name=`"$key`"
		value=`"$value`" />
"@
	}
	Write-Output "</service>"
}

$proxyServer = 'http://proxy-server.:3271/'
function Invoke-WebRequest-Wrap() {
	param (
		[Parameter(Mandatory)][string]$Uri,
		[Parameter()][string]$Out
	)

	#  -Authentication Basic -Credential ""
	Invoke-WebRequest -UserAgent 'GongT/wireguard-config-distribute' -Proxy $proxyServer -Uri $Uri -Out $Out
}

sc.exe query wg_normal | Out-Null
if ($lastexitcode -eq 0) {
	Write-Host "删除旧服务……"
	sc.exe stop wg_normal | Out-Null
	sc.exe delete wg_normal
	sc.exe query wg_normal | Out-Null
	if ($lastexitcode -eq 0) {
		Write-Error "删除失败"
	}
}

if ($env:OneDriveConsumer) {
	$Root = "$env:OneDriveConsumer/Software/WireguardConfig"
} elseif ($env:OneDrive) {
	$Root = "$env:OneDrive/Software/WireguardConfig"
} else {
	Write-Error "木有找到 OneDrive 路径"
}
$configFile = "$Root/$env:COMPUTERNAME.conf"
if (-Not (Test-Path $configFile)) {
	Write-Error "配置文件不存在：$configFile"
}
Write-Host "使用的配置文件：$configFile"

$hashTable = Get-Content -Raw -Path $configFile | ConvertFrom-StringData
foreach ($key in $hashTable.Keys) {
	$value = $hashTable.$key
	Set-Item env:$key $value 
}
if (!$env:WIREGUARD_GROUP) {
	Write-Error "缺少必须的设置：WIREGUARD_GROUP"
}

$distFolder = "C:/Program Files/WireGuard"
if (-Not (Test-Path "$distFolder/wireguard.exe")) {
	Write-Warning "WireGuard可执行文件不存在，别忘了安装！"
	if (-Not (Test-Path $distFolder)) {
		New-Item -ItemType Directory $distFolder | Out-Null
	}
}

$winswBinary = (downloadGithubRelease -repo 'winsw/winsw' -getfile 'WinSW.NET461.exe' -saveas 'winsw.exe' -distfolder $distFolder)
Write-Host "winsw版本：" -NoNewline
& $winswBinary --version
if ($lastexitcode -ne 0) { Write-Error "程序无法运行！($lastexitcode)" }

$clientBinary = "$(Split-Path -Parent $winswBinary)/wireguard-config-distribute-client.exe"
$downloadUrl = buildGithubReleaseUrl -Repo 'GongT/wireguard-config-distribute' -TagName 'latest' -GetFile 'client.exe'
# $clientBinary = (downloadGithubRelease -repo 'GongT/wireguard-config-distribute' -getfile 'client.exe' -saveas 'wireguard-config-distribute-client.exe' -distfolder $distFolder)
# Write-Host "客户端版本：" -NoNewline
# & $clientBinary /version
# if ($lastexitcode -ne 0)  {
# 	Write-Error "程序无法运行！"
# }

$serviceConfigFile = "$distFolder/$($env:WIREGUARD_GROUP).xml"
createConfig | Out-File -Encoding utf8 -FilePath $serviceConfigFile

& $winswBinary test $serviceConfigFile
if ($lastexitcode -ne 0) {
	Write-Host "安装Windows服务……"
	& $winswBinary install $serviceConfigFile
	if ($lastexitcode -ne 0) { 
		Write-Error "安装失败，返回：$lastexitcode"
	}
} else {
	& $winswBinary refresh $serviceConfigFile
}


Write-Host "启动……"
& $winswBinary start $serviceConfigFile
if ($lastexitcode -ne 0) { 
	Write-Error "启动失败，返回：$lastexitcode"
}

Write-Host "完成！"
