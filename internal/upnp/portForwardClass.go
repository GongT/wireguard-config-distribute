// +build !android

package upnp

import (
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

type UPnPPortForwarder struct {
	ch       chan uint16
	OnChange <-chan uint16

	gatewayIP      net.IP
	client         *natpmp.Client
	keepAliveTimer *time.Ticker

	currentExternal uint16

	listenPort int
	wantPort   int

	pDispose func()
}

type IOptions interface {
	GetPublicPort() uint16
	GetListenPort() uint16
}

func NewAutoForward(opts IOptions) (*UPnPPortForwarder, error) {
	ret := &UPnPPortForwarder{
		wantPort:   int(opts.GetPublicPort()),
		listenPort: int(opts.GetListenPort()),
	}

	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return nil, err
	}

	ret.gatewayIP = gatewayIP

	client := natpmp.NewClientWithTimeout(gatewayIP, 10*time.Second)
	ret.client = client

	ret.ch = make(chan uint16)
	ret.OnChange = ret.ch

	ret.pDispose = tools.WaitExit(func(int) {
		ret.Stop()
	})

	return ret, nil
}
func (p *UPnPPortForwarder) Stop() {
	tools.Debug("[UPnP] Stop()")
	p.pDispose()
	if p.keepAliveTimer != nil {
		p.keepAliveTimer.Stop()
		close(p.ch)
	}
	p.client.AddPortMapping("udp", p.listenPort, p.wantPort, 0)
}

func (p *UPnPPortForwarder) Start() (uint16, error) {
	fw, err := p.client.AddPortMapping("udp", p.listenPort, p.wantPort, 2*60)
	if err != nil {
		tools.Error("[UPnP] AddPortMapping() fail: %v | will not try again", err)
		return 0, err
	}
	p.currentExternal = fw.MappedExternalPort

	tools.Debug("[UPnP] AddPortMapping(): %v", p.currentExternal)

	p.keepAliveTimer = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-p.keepAliveTimer.C:
				p.tick()
			}
		}
	}()

	return fw.MappedExternalPort, nil
}

func (p *UPnPPortForwarder) tick() {
	fw, err := p.client.AddPortMapping("udp", p.listenPort, p.wantPort, 2*60)
	if err != nil {
		tools.Error("[UPnP] AddPortMapping() fail: %v", err)
		return
	}

	if p.currentExternal != fw.MappedExternalPort {
		tools.Error("[UPnP] AddPortMapping() open port change from %v to %v.", p.currentExternal, fw.MappedExternalPort)
		p.currentExternal = fw.MappedExternalPort
		p.ch <- fw.MappedExternalPort
	}
}
