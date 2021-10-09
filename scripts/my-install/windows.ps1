#!/usr/bin/env pwsh

$proxyServer = 'http://proxy-server.:3271/'
$winswBinaryDownloadName = 'WinSW-net461.exe'

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
	$root = "$env:tmp/__wgc_install/"
	if (-Not (Test-Path "$root/scripts/my-install")) {
		New-Item -Path "$root/scripts/my-install" -ItemType Directory | Out-Null
	}
	if (-Not (Test-Path "$root/scripts/services")) {
		New-Item -Path "$root/scripts/services" -ItemType Directory | Out-Null
	}
	cp $PSScriptRoot/windows.ps1 "$root/scripts/my-install"
	cp $PSScriptRoot/../services/auto-update.ps1 "$root/scripts/services"
	Start-Process -Verb RunAs pwsh -ArgumentList '-NoExit', "$root/scripts/my-install/windows.ps1"
	exit
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

function getLocalVersion() {
	$binaryFile = "$distFolder/$getFile"
	if (Test-Path -Path $binaryFile  ) {
		$v = & $binaryFile --version 2>$null | Select-String -Pattern "Hash:(.+)"
		return $v.Matches.Groups[1].Value.Trim()
	} else {
		return ""
	}
}

function detectVersionChange() {	
	param (
		[Parameter(Mandatory)][string]$Repo,
		[Parameter(Mandatory)][string]$versionFile,
		[Parameter(Mandatory)][scriptblock]$callback
	)

	$releaseDataUrl = "https://api.github.com/repos/$Repo/releases?page=1&per_page=1"
	
	Write-Host "检查 $Repo 版本……"

	if (Test-Path -Path $versionFile) {
		Write-Host -ForegroundColor Gray "    记录文件： $versionFile"
		$versionLocal = Get-Content -Encoding utf8 $versionFile
	} else {
		Write-Host -ForegroundColor Gray "    记录文件： 不存在"
		$versionLocal = ''
	}

	Write-Host -ForegroundColor Gray "    来源： $releaseDataUrl"
	$releaseData = (Invoke-WebRequest-Wrap -Uri $releaseDataUrl | ConvertFrom-Json)[0]

	$downloadUrl = buildGithubReleaseUrl -Repo $Repo -TagName $releaseData.tag_name -GetFile $GetFile
	if ($versionLocal -eq $releaseData.target_commitish) {
		Write-Host -ForegroundColor Gray "    ~ 没有更新"
		return Invoke-Command $callback -ArgumentList $false, $downloadUrl
	} else {
		Write-Host -ForegroundColor Gray "    -> 有更新: $versionLocal → 远程：$($releaseData.target_commitish)"
		$ret = Invoke-Command $callback -ArgumentList $true, $downloadUrl
		Set-Content -Encoding utf8 -Path $versionFile -Value $releaseData.target_commitish
		return $ret
	}
}

function downloadGithubRelease() {
	param (
		[Parameter(Mandatory)][string]$Repo,
		[Parameter(Mandatory)][string]$GetFile,
		[Parameter()][string]$SaveAs = $GetFile,
		[Parameter(Mandatory)][string]$DistFolder
	)

	$versionFile = "$DistFolder/$SaveAs.version.txt"

	$binaryFile = detectVersionChange -Repo $Repo -versionFile $versionFile -callback {
		param($change, $downloadUrl)
		$binaryFile = "$DistFolder/$SaveAs"
		$downloadFile = "$DistFolder/$SaveAs.update"

		if (-Not ($change)) {
			Write-Host " * 已是最新版本"
			Write-Host -ForegroundColor Gray "    文件:   $binaryFile"
			return $binaryFile
		}
	
		Remove-Item $binaryFile -ErrorAction SilentlyContinue | Out-Null
		Remove-Item $downloadFile -ErrorAction SilentlyContinue | Out-Null

		# $releaseData.assets
		Write-Host " * 有更新，开始下载："
		Write-Host -ForegroundColor Gray "    远程: $downloadUrl"
		Write-Host -ForegroundColor Gray "    本地:   $downloadFile"
		Invoke-WebRequest-Wrap -Uri $downloadUrl -OutFile $downloadFile

		copyFileIf $downloadFile $binaryFile | Out-Null

		return $binaryFile
	}

	return $binaryFile
}

function copyFileIf() {
	param($from, $to)
	
	$TempFile = New-TemporaryFile
	Copy-Item -Path $from -Destination $TempFile | Out-Null
	$from = $TempFile	
	
	$MOVEFILE_REPLACE_EXISTING = 0x1
	$MOVEFILE_DELAY_UNTIL_REBOOT = 0x4
	$MOVEFILE_WRITE_THROUGH = 0x8
		
	### https://gist.github.com/marnix/7565364
	### https://stackoverflow.com/questions/24391367/dllimport-in-powershell-for-accessing-c-style-32-bit-api-using-relative-path
	$signature = @'
		[DllImport("kernel32.dll", SetLastError = true, CharSet = CharSet.Unicode)]
		public static extern bool MoveFileEx(string lpExistingFileName, string lpNewFileName, UInt32 dwFlags);
'@
	Add-Type -MemberDefinition $signature -Name "MoveFile" -Namespace Win32Function

	if (Move-Item -Path $from -Destination $to -Force) {
		return $true
	}

	Write-Host "文件被占用，将于下次重启时更新。"
	$deleteResult = [Win32Function.MoveFile]::MoveFileEx($from, $to, $MOVEFILE_REPLACE_EXISTING + $MOVEFILE_DELAY_UNTIL_REBOOT)

	if ($deleteResult -eq $false) {
		throw (New-Object ComponentModel.Win32Exception) # calls GetLastError
	}
	return $false
}

function createConfig() {
	Write-Output @"
<service>
	<id>wg_$($env:WIREGUARD_GROUP)</id>
	<name>Wireguard 自动配置[组：$($env:WIREGUARD_GROUP)]</name>
	<executable>$clientBinary</executable>
	<delayedAutoStart>true</delayedAutoStart>
	<stoptimeout>20sec</stoptimeout>
	<onfailure
		action="restart"
		delay="20sec" />
	<serviceaccount>
		<username>LocalSystem</username>
		<allowservicelogon>true</allowservicelogon>
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

function stringifyFunction() {
	[CmdletBinding()]
	param (
		[Parameter(Mandatory)][string]$Fn
	)
	$ret = "function $Fn() {`n"
	$ret += (Get-Command $Fn).Definition
	$ret += "`n}"
	return $ret
}

function stopAllService() {
	if (-Not (Test-Path "$distFolder/winsw.exe")) {
		return
	}
	$serviceConfigList = Get-ChildItem -Path $distFolder -Depth 1 -Filter '*.xml'
	foreach ($item in $serviceConfigList ) {
		Write-Host "    停止服务：$item"
		& "$distFolder/winsw.exe" stop $item
	}
}

function startAllService() {
	$serviceConfigList = Get-ChildItem -Path $distFolder -Depth 1 -Filter '*.xml'
	foreach ($item in $serviceConfigList ) {
		Write-Host "    启动服务：$item"
		& "$distFolder/winsw.exe" start $item
	}
}

function createUpdateSchedule() {
	$taskPath = "\GongT\"
	$taskName = "wireguard-config-client-auto-update"
	try {
		Unregister-ScheduledTask -Confirm:$false -TaskName $taskName -TaskPath $taskPath
	} catch {
	}

	$acl = Get-Acl $distFolder
	$AccessRule = New-Object -TypeName System.Security.AccessControl.FileSystemAccessRule `
		-ArgumentList "NT AUTHORITY\NetworkService", "FullControl", "Allow" 
	$acl.SetAccessRule($AccessRule)
	Set-Acl $distFolder -AclObject $acl

	# $OnBoot = New-ScheduledTaskTrigger -AtStartup
	$Timer = New-ScheduledTaskTrigger -Daily -At "00:00" 

	# $currentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent().Name
	# $User = New-ScheduledTaskPrincipal -UserId $currentUser -LogonType S4U -RunLevel Highest
	$User = New-ScheduledTaskPrincipal -UserId "NT AUTHORITY\NetworkService" -LogonType ServiceAccount -RunLevel Highest
	# $User = New-ScheduledTaskPrincipal -GroupId "BUILTIN\Administrators" -RunLevel Highest

	$powershell = "$PSHOME\pwsh.exe"
	$RunScript = New-ScheduledTaskAction -Execute $powershell -Argument "-File auto-update.ps1" -WorkingDirectory $distFolder -Id "auto-update"

	$Settings = New-ScheduledTaskSettingsSet `
		-StartWhenAvailable `
		-AllowStartIfOnBatteries `
		-Compatibility Win8 `
		-RunOnlyIfNetworkAvailable `
		-ExecutionTimeLimit (New-TimeSpan -Minutes 5) `
		-MultipleInstances IgnoreNew `
		-DontStopIfGoingOnBatteries `
		-RestartCount 3 `
		-RestartInterval (New-TimeSpan -Minutes 1)

	$TaskInstance = New-ScheduledTask -Description "Wireguard配置工具自动更新程序" -Settings $Settings -Trigger $Timer -Principal $User -Action $RunScript

	Register-ScheduledTask -TaskName $taskName -TaskPath $taskPath -InputObject $TaskInstance | Out-Null
}

function Invoke-WebRequest-Wrap() {
	param (
		[Parameter(Mandatory)][Uri]$Uri,
		[Parameter()][string]$OutFile
	)

	$param = @{}
	if ($OutFile) {
		$param.OutFile = $OutFile
		$param.Resume = $true
	}
	if ($proxyServer) {
		$param.Proxy = $proxyServer
	}

	$tokenFile = "$Root/.github-token"
	if (Test-Path -Path $tokenFile -PathType Leaf) {
		$token = Get-Content -Encoding utf8 $tokenFile
		$param.Headers = @{Authorization = "token $token" }
	}
	try {
		$response = Invoke-WebRequest @param `
			-MaximumRetryCount 10 `
			-RetryIntervalSec 5 `
			-UserAgent 'GongT/wireguard-config-distribute' `
			-Uri $Uri
		return $response
	} catch {
		Write-Host $_.Exception
		exit 1
	}
}

# sc.exe query wg_normal | Out-Null
# if ($lastexitcode -eq 0) {
# 	Write-Host "删除旧服务……"
# 	sc.exe stop wg_normal | Out-Null
# 	sc.exe delete wg_normal
# 	sc.exe query wg_normal | Out-Null
# 	if ($lastexitcode -ne 0) {
# 		$le = $lastexitcode
# 		sc.exe query wg_normal
# 		Write-Error "删除失败，返回：$le"
# 	}
# }

if ($env:OneDriveConsumer) {
	$Root = "$env:OneDriveConsumer/AppData/WireguardConfig"
} elseif ($env:OneDrive) {
	$Root = "$env:OneDrive/AppData/WireguardConfig"
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

stopAllService

$winswBinary = downloadGithubRelease -repo 'winsw/winsw' -getfile $winswBinaryDownloadName -saveas 'winsw.exe' -distfolder $distFolder
Write-Host "winsw版本：$(& $winswBinary --version)"
if ($lastexitcode -ne 0) { Write-Error "程序无法运行！($lastexitcode)" }

$clientBinary = downloadGithubRelease -repo 'GongT/wireguard-config-distribute' -getfile 'client.exe' -distfolder $distFolder
Write-Host "客户端版本：$(& $clientBinary --version)"
if ($lastexitcode -ne 0) {
	Write-Error "程序无法运行！"
}

$serviceConfigFile = "$distFolder/$($env:WIREGUARD_GROUP).xml"
createConfig | Out-File -Encoding utf8 -FilePath $serviceConfigFile

Write-Host "安装任务计划……"
Copy-Item "$PSScriptRoot/../services/auto-update.ps1" $distFolder | Out-Null
Set-Content -Path "$distFolder/lib.ps1" -Value @(
	(stringifyFunction Invoke-WebRequest-Wrap),
	(stringifyFunction buildGithubReleaseUrl),
	(stringifyFunction detectVersionChange),
	(stringifyFunction copyFileIf),
	(stringifyFunction stopAllService),
	(stringifyFunction startAllService)
)
createUpdateSchedule

sc.exe query wg_normal | Out-Null
if ($lastexitcode -eq 0) {
	Write-Host "更新Windows服务定义……"
	& $winswBinary refresh $serviceConfigFile
} else {
	Write-Host "安装Windows服务……"
	& $winswBinary install $serviceConfigFile
	if ($lastexitcode -ne 0) { 
		Write-Error "安装失败，返回：$lastexitcode"
	}
}


startAllService

Write-Host "完成！"
