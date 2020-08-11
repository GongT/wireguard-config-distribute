package wireguardControl

import (
	"fmt"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
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
		tools.Debug("Create extended config file: %v", exCfg)
		if err := saveBuffersTo(exCfg, wc.creatConfigHeader(true), wc.creatConfigBody()); err != nil {
			return fmt.Errorf("failed write file [%s]: %v", exCfg, err)
		}
	}
	return nil
}
