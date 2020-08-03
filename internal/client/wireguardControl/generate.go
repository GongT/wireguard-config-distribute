package wireguardControl

import (
	"bytes"
	"fmt"
)

func (wc *WireguardControl) creatConfig() []byte {
	result := bytes.NewBuffer(make([]byte, 0, 2048))

	appendLine := func(line string, args ...interface{}) {
		result.WriteString(fmt.Sprintf(line, args...))
		result.WriteByte('\n')
	}

	appendLine("[Interface]")
	appendLine("# Name = %s (%s)", wc.interfaceTitle)
	appendLine("Address = %s/32", wc.givenAddress)
	appendLine("ListenPort = %d", wc.interfaceListenPort)
	appendLine("PrivateKey = %s", wc.privateKey)
	// appendLine("DNS = 1.1.1.1,8.8.8.8")
	// appendLine("Table = 12345")
	if wc.interfaceMTU > 0 {
		appendLine("MTU = %d", wc.interfaceMTU)
	}
	appendLine("")

	for _, peer := range wc.peers {
		appendLine("[Peer]")
		appendLine("# Name = %s", peer.comment)
		if wc.subnet > 0 {
			appendLine("AllowedIPs = %s/%d", peer.privateIp, wc.subnet)
		} else {
			appendLine("AllowedIPs = %s/32", peer.privateIp)
		}
		appendLine("Endpoint = %s:%d", peer.ip, peer.port)
		appendLine("PublicKey = %s", peer.publicKey)
		if len(peer.presharedKey) > 0 {
			appendLine("PresharedKey = %s", peer.presharedKey)
		}
		if peer.keepAlive > 0 {
			appendLine("PersistentKeepalive = %d", peer.keepAlive)
		}
		appendLine("")
	}

	return result.Bytes()
}
