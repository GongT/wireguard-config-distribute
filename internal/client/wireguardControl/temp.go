package wireguardControl

import (
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func tryTempDir(dir string) (string, bool) {
	if dir == "" {
		return "", false
	}
	ret := filepath.Join(dir, "wireguard")
	err := os.MkdirAll(ret, os.FileMode(0755))
	if err != nil {
		tools.Error("failed create dir [%s]: %s", ret, err.Error())
		return "", false
	}

	return ret, true
}

var TempDir string

func init() {
	if ret, ok := tryTempDir(os.Getenv("XDG_CACHE_HOME")); ok {
		TempDir = ret
		return
	}
	if ret, ok := tryTempDir(os.Getenv("CACHE_DIRECTORY")); ok {
		TempDir = ret
		return
	}
	if ret, ok := tryTempDir(os.TempDir()); ok {
		TempDir = ret
		return
	}
	tools.Die("failed find location to store temp file(s)")
}
