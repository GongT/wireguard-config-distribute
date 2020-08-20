package interfaceState

import (
	"fmt"
	"log"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/vishvananda/netlink"
)

type InterfaceOptions interface {
	GetNetwork() string
	GetAddress() string
	GetMtu() uint16
}

type nativeState struct {
}

func (is *interfaceState) init() {
	_, err := netlink.LinkList()
	if err != nil {
		tools.Die("failed to call netlink api: %s", err.Error())
	}
}

func (is *interfaceState) DeleteInterface() error {
	if link, err := getLink(is.ifname); link != nil {
		err := netlink.LinkDel(link)
		if err != nil {
			return fmt.Errorf("failed delete network interface: %v", err)
		}
		is.network = ""
		is.mtu = 0
	} else if err != nil {
		tools.Error("failed delete(step get) network interface: %v", err)
	}
	return nil
}

func (is *interfaceState) CreateOrUpdateInterface(options InterfaceOptions) error {
	log.Println("Creating network interface!")
	if link, err := getLink(is.ifname); err != nil {
		return fmt.Errorf("error when get network interface: %v", err)
	} else if link == nil {
		return is.create(options)
	} else {
		changed := diffState(is, options)
		if changed.network {
			if err := is.set_ip(link, options); err != nil {
				return err
			}
		}
		if changed.mtu {
			tools.Debug("interface configure has changed: MTU: %v -> %v", is.mtu, options.GetMtu())
			err := netlink.LinkSetMTU(link, int(options.GetMtu()))
			if err != nil {
				return fmt.Errorf("failed set network interface MTU: %s", err)
			}
		}
		changed.commit()
	}
	return nil
}

func (is *interfaceState) create(options InterfaceOptions) error {
	la := netlink.NewLinkAttrs()
	if options.GetMtu() > 0 {
		la.MTU = int(options.GetMtu())
	}
	la.Name = is.ifname
	link := &netlink.GenericLink{
		LinkAttrs: la,
		LinkType:  "wireguard",
	}

	if err := netlink.LinkAdd(link); err != nil {
		return fmt.Errorf("failed create network interface: %s", err)
	}
	is.mtu = options.GetMtu()

	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed set interface up: %s", err)
	}

	return is.set_ip(link, options)
}

func (is *interfaceState) set_ip(link netlink.Link, options InterfaceOptions) error {
	address := options.GetAddress()

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return fmt.Errorf("failed parse address [%s]: %v", address, err)
	}

	tools.Debug("<if> ip add %v", addr.String())
	if err := netlink.AddrReplace(link, addr); err != nil {
		return fmt.Errorf("failed replace interface address: %v", err)
	}

	// XX dev wg_YY proto static scope link
	networkStr := options.GetNetwork()
	_, network, err := net.ParseCIDR(networkStr)
	if err != nil {
		return fmt.Errorf("failed parse network [%s]: %v", networkStr, err)
	}

	tools.Debug("<if> route add %v", network.String())
	netlink.RouteReplace(&netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       network,
		Src:       addr.IP,
		Table:     254,
	})

	is.network = networkStr
	is.address = address

	return nil
}

func getLink(name string) (netlink.Link, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			return nil, nil
		}
		return nil, fmt.Errorf("failed LinkByName: %s", err.Error())
	}
	if link.Type() != "wireguard" {
		return nil, fmt.Errorf("link `%s' is typeof `%s' but required `wireguard'", name, link.Type())
	}
	return link, nil
}
