package wireguardControl

import (
	"golang.zx2c4.com/wireguard/wgctrl"
)

type InterfaceControl struct {
	client *wgctrl.Client
	opts   InterfaceOptions
}

type InterfaceOptions interface {
	GetListenPort() uint16
	GetInterfaceName() string
	GetMTU() uint16
}

func CreateInterfaceControl(opts InterfaceOptions) *InterfaceControl {
	// client := wgctrl.New()
	return nil
}
