package wireguardControl

import (
	"path/filepath"
)

func (wc *WireguardControl) GetConfigFile() (string, error) {
	if err := wc.createExtendConfigFile(); err != nil {
		return "", err
	}
	return filepath.Join(TempDir, wc.interfaceName+".conf"), nil
}
