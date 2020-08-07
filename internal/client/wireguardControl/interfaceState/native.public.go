package interfaceState

import "github.com/gongt/wireguard-config-distribute/internal/tools"

type publicIfOptions interface {
	GetNetwork() string
	GetMtu() int
}

type changed struct {
	network bool
	mtu     bool

	opts publicIfOptions
	is   *interfaceState
}

func diffState(prevStat *interfaceState, newStat publicIfOptions) (r changed) {
	r.is = prevStat
	r.opts = newStat
	if prevStat.network != newStat.GetNetwork() {
		tools.Debug("interface configure has changed: address: [%v] -> [%v]", prevStat.network, newStat.GetNetwork())
		r.network = true
	}
	if prevStat.mtu != newStat.GetMtu() {
		tools.Debug("interface configure has changed: MTU: %v -> %v", prevStat.mtu, newStat.GetMtu())
		r.mtu = true
	}
	return
}

func (c *changed) commit() {
	if c.mtu {
		c.is.mtu = c.opts.GetMtu()
	}
	if c.network {
		c.is.network = c.opts.GetNetwork()
	}
}
