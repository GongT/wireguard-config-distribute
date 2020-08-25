package clientType

import "github.com/gongt/wireguard-config-distribute/internal/protocol"

type PeerData struct {
	*protocol.Peers_Peer

	selectedAddress string
	selectedPort    uint16
}

type PeerDataList []*PeerData

func (pd PeerData) GetSelectedAddress() string {
	return pd.selectedAddress
}
func (pd PeerData) GetSelectedPort() uint16 {
	return pd.selectedPort
}
func (pd *PeerData) ChangeTo(ip string, port uint16) {
	pd.selectedAddress = ip
	pd.selectedPort = port
}

func WrapList(list []*protocol.Peers_Peer, ipv4Only bool) PeerDataList {
	ret := make(PeerDataList, 0, len(list))
	for _, client := range list {
		peer := client.GetPeer()
		var selectedIp string
		if peer.GetSameNetwork() {
			selectedIp = peer.GetAddress()[0]
		} else {
			selectedIp = selectIp(peer.GetAddress(), ipv4Only)
		}

		ret = append(ret, &PeerData{
			Peers_Peer:      client,
			selectedAddress: selectedIp,
			selectedPort:    uint16(peer.GetPort()),
		})
	}
	return ret
}
