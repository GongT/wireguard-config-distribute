package wireguard

import (
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type KeyPair struct {
	Public  string
	Private string
}

func ParseKey(private string) (*KeyPair, error) {
	key, err := wgtypes.ParseKey(private)
	if err != nil {
		return nil, fmt.Errorf("Failed parse Wireguard private key: %s: %v", private, err)
	}
	return &KeyPair{
		Private: key.PublicKey().String(),
		Public:  key.PublicKey().String(),
	}, nil
}

func AllocateKeyPair() (*KeyPair, error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Private: key.PublicKey().String(),
		Public:  key.PublicKey().String(),
	}, nil
}
