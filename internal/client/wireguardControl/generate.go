package wireguardControl

import (
	"fmt"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (wc *WireguardControl) creatConfigHeader(extendedSyntax bool) []byte {
	result := newBuffer(extendedSyntax)

	result.appendLine("[Interface]")
	result.appendLine("# Name = %s", wc.interfaceTitle)
	result.appendLineExtened("Address = %s/32", wc.givenAddress)
	result.appendLine("ListenPort = %d", wc.interfaceListenPort)
	result.appendLine("PrivateKey = %s", wc.privateKey)
	// appendLineExtened("DNS = 1.1.1.1,8.8.8.8")
	// appendLineExtened("Table = 12345")
	if wc.lowestMtu > 0 {
		result.appendLineExtened("MTU = %d", wc.lowestMtu)
	}
	// PreUp, PostUp, PreDown, PostDown
	// SaveConfig?
	result.appendLine("")

	return result.Bytes()
}

func (wc *WireguardControl) creatConfigBody() []byte {
	result := newBuffer(false)

	for _, peer := range wc.peers {
		result.appendLine("[Peer]")
		result.appendLine("# Name = %s", peer.comment)
		// if wc.subnet > 0 {
		// 	result.appendLine("AllowedIPs = %s/%d", peer.privateIp, wc.subnet)
		// } else {
		result.appendLine("AllowedIPs = %s/32", peer.privateIp)
		// }
		if len(peer.ip) > 0 {
			if strings.Contains(peer.ip, ":") {
				result.appendLine("Endpoint = [%s]:%d", peer.ip, peer.port)
			} else {
				result.appendLine("Endpoint = %s:%d", peer.ip, peer.port)
			}
		} else {
			result.appendLine("# Endpoint is not public accessable")
		}
		result.appendLine("PublicKey = %s", peer.publicKey)
		if len(peer.presharedKey) > 0 {
			result.appendLine("PresharedKey = %s", peer.presharedKey)
		}
		if peer.keepAlive > 0 {
			result.appendLine("PersistentKeepalive = %d", peer.keepAlive)
		}
		result.appendLine("")
	}

	return result.Bytes()
}

func (wc *WireguardControl) createConfigFile() error {
	wc.extendedConfigCreated = false
	tools.Debug("Create native config file: %v", wc.configFile)
	if err := saveBuffersTo(wc.configFile, wc.creatConfigHeader(false), wc.creatConfigBody()); err != nil {
		return fmt.Errorf("failed write file [%s]: %v", wc.configFile, err)
	}
	return nil
}
