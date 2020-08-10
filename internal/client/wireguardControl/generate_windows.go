package wireguardControl

import (
	"fmt"
	"path/filepath"
)

func (wc *WireguardControl) GetConfigFile() (string, error) {
	if err := wc.createExtendConfigFile(); err != nil {
		return "", err
	}
	return filepath.Join(TempDir, wc.interfaceName+".conf"), nil
}

func (wc *WireguardControl) createExtendConfigFile() error {
	if !wc.extendedConfigCreated {
		exCfg := filepath.Join(TempDir, wc.interfaceName+".conf")
		if err := saveBuffersTo(exCfg, wc.creatConfigHeader(false), wc.creatConfigBody()); err != nil {
			return fmt.Errorf("failed write file [%s]: %v", exCfg, err)
		}
	}
	return nil
}
