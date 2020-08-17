#!/usr/bin/env bash

set -Eeuo pipefail

cd "$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
cd ../..


# ./build/tools/make-standalone-toolchain.sh --platform=android-21 --install-dir=/data/DevelopmentRoot/Android/NDK --arch=arm64
export CC="/data/DevelopmentRoot/Android/NDK/bin/armv7a-linux-androideabi28-clang"

pwsh scripts/build.ps1 android
adb push dist/client.android /data/local/tmp/
