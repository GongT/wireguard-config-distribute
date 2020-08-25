package peerStatus

import "fmt"

func (peers *PeersManager) Dump() string {
	defer peers.m.Lock("Dump")()

	ret := "== Peers Mamager Status ==\n"

	ret += fmt.Sprintf("GUID Map: (next=%v)\n", peers.guid)
	for id, guid := range peers.guidMap {
		ret += fmt.Sprintf("  * %v -> %v\n", id, guid)
	}
	ret += "\n"

	for vpn, peersList := range peers.mapper {
		ret += fmt.Sprintf("Peers List: %v (count=%v)\n", vpn.Serialize(), len(peersList))
		for guid, peer := range peersList {
			ret += fmt.Sprintf(`> %v - %v
    MachineId: %v
    VpnName: %v
    Hostname: %v
    PublicKey: %v
    VpnIp: %v
    MTU: %v
    HostsLine: %v
    NetworkId: %v
    ExternalIp: %v
    ExternalPort: %v
    InternalIp: %v
    InternalPort: %v
    lastKeepAlive: %v
`,
				guid, peer.Title, peer.MachineId, peer.VpnId.Serialize(), peer.Hostname, peer.PublicKey,
				peer.VpnIp, peer.MTU, peer.HostsLine, peer.WorkgroupId, peer.ExternalIp,
				peer.ExternalPort, peer.InternalIp, peer.InternalPort, peer.lastKeepAlive,
			)
		}
	}
	return ret
}
