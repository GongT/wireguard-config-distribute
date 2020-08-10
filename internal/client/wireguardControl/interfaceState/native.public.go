package interfaceState

import "github.com/gongt/wireguard-config-distribute/internal/tools"

type publicIfOptions interface {
	GetNetwork() string
	GetAddress() string
	GetMtu() int
}

type changed struct {
	network bool
	address bool
	mtu     bool

	opts publicIfOptions
	is   *interfaceState
}

func diffState(prevStat *interfaceState, newStat publicIfOptions) (r changed) {
	r.is = prevStat
	r.opts = newStat
	if prevStat.address != newStat.GetAddress() {
		tools.Debug("interface configure has changed: address: [%v] -> [%v]", prevStat.network, newStat.GetAddress())
		r.address = true
	}
	if prevStat.mtu != newStat.GetMtu() {
		tools.Debug("interface configure has changed: MTU: %v -> %v", prevStat.mtu, newStat.GetMtu())
		r.mtu = true
	}
	if prevStat.network != newStat.GetNetwork() {
		tools.Debug("interface configure has changed: network: %v -> %v", prevStat.network, newStat.GetNetwork())
		r.network = true
	}
	return
}

func (c *changed) commit() {
	if c.mtu {
		c.is.mtu = c.opts.GetMtu()
	}
	if c.address {
		c.is.address = c.opts.GetAddress()
	}
	if c.network {
		c.is.network = c.opts.GetNetwork()
	}
}
