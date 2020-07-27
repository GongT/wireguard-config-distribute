package wireguard

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func GenerateKeyPair() (pub string, pri string, err error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return
	}
	pri = key.String()
	pub = key.PublicKey().String()

	return
}
