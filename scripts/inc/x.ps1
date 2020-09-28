# Write-Output "!!! global:execWithContainer = $($global:execWithContainer)$($global:CID)"

function hostx() {
	$cmd, $exargs = $args
	Write-Host -Separator " " -ForegroundColor DarkGray " +" $cmd $exargs
	& $cmd @exargs
}

function podmanx() {
	hostx podman run --rm `
		--http-proxy `
		"--workdir=/app" `
		"--volume=$(Get-Location):/app" `
		"--volume=$(go env GOCACHE):/GOCACHE" `
		$global:CID `
		@args
}

function  SetExecuteMethod() {
	param (
		[switch]$container
	)
	if ($container) {
		Write-Output "do everything inside a container"

		$global:execWithContainer = $true 

		$global:CID = $(podman inspect --type=image '--format={{.Id}}' gongt/wg-config-build)
		if ($? -ne $true) {
			if (-Not $env:GOPATH) {
				if ( Test-CommandExists go ) {
					$env:GOPATH = $(go env GOPATH)
				}
				if (-Not $env:GOPATH) {
					$env:GOPATH = New-Item -Force -ItemType "directory" -Path $env:SYSTEM_COMMON_CACHE -Name gopath
					Write-Output "using dummy GOPATH=${env:GOPATH}"
				} else {
					Write-Output "using go env GOPATH=${env:GOPATH}"
				}
			} else {
				Write-Output "using system GOPATH=${env:GOPATH}"
			}
			$cache = @("--volume=$(go env GOCACHE):/GOCACHE")
			if ( Test-Path $env:SYSTEM_COMMON_CACHE/apk ) {
				$cache += "--volume=$env:SYSTEM_COMMON_CACHE/apk:/etc/apk/cache"
			}
			Write-Host -ForegroundColor Gray " + podman build $cache --file - --tag gongt/wg-config-build"
			Write-Output "
				FROM gongt/alpine-cn:edge
				ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
				ENV GOBIN=/usr/bin GOPATH=/root/go GOPROXY=$(go env GOPROXY) GOCACHE=/GOCACHE GO111MODULE=auto
				RUN cd / \
				 && set -x \
				 && apk add -U go protoc git protobuf-dev \
				 && go get -v -u github.com/GongT/go-generate-struct-interface/cmd/go-generate-struct-interface github.com/golang/protobuf/protoc-gen-go \
				 && command -v protoc-gen-go || ( find / -name protoc-gen-go ; exit 1) \
				 && command -v go-generate-struct-interface || ( find / -name go-generate-struct-interface ; exit 1)
			" | podman build @cache --file - --tag "gongt/wg-config-build"
			if ($? -eq $false) {
				Write-Error "Failed create image for build"
				exit 1
			}

			$global:CID = $(podman inspect --type=image '--format={{.Id}}' gongt/wg-config-build)
		}
		
		# Write-Output "execWithContainer: $($global:execWithContainer), CID: $($global:CID)"
	}
}
function x() {
	if ($global:execWithContainer -eq $true) {
		podmanx @args
	} else {
		hostx @args
	}
}
