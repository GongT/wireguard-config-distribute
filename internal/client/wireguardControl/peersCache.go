package wireguardControl

import "net"

type peerData struct {
	publicKey    string
	presharedKey string
	ip           net.IP
	port         uint
	keepAlive    uint
	privateIp    net.IP
}

type PeersCache struct {
	status map[string]peerData
}

func NewPeersCache() *PeersCache {
	return nil
}
