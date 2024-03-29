package detect_ip

import (
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type Options interface {
	GetPublicIp() []string
	GetGateway() bool
	GetIpUpnpDisable() bool
	GetIpApi() string
	GetNoPublicNetwork() bool
}

type Detect struct {
	disabled bool

	readInterface bool
	useUPnP       bool
	api           string

	manualSet []string
	lastGet   []string
}

func NewDetect(options Options) *Detect {
	ret := Detect{
		disabled:      options.GetNoPublicNetwork(),
		readInterface: options.GetGateway(),
		useUPnP:       !options.GetIpUpnpDisable(),
		api:           options.GetIpApi(),
		manualSet:     options.GetPublicIp(),
	}

	return &ret
}

func (d *Detect) GetLast() []string {
	return d.lastGet
}

func (d *Detect) GetDisabled() bool {
	return d.disabled
}

func (d *Detect) Execute() {
	if d.disabled {
		tools.Debug("Skip public ip detect: no public access")
		d.lastGet = make([]string, 0)
		return
	}

	ret := make([]string, 0, len(d.lastGet)+10)

	for _, ip := range d.manualSet {
		ret = append(ret, ip)
	}

	gotIpv4 := false

	tools.Debug("get ip address from local interfaces:")
	for _, ip := range ListAllLocalNetworkIp() {
		if tools.IsIPv4(ip) {
			if d.readInterface && IsPublicIp(ip) {
				gotIpv4 = true
				ret = append(ret, ip.String())
				tools.Debug("  -> ipv4: %v", ip.String())
			} else {
				tools.Debug("  x> %v", ip.String())
			}
		} else {
			tools.Debug("  x> %v", ip.String())
		}
	}

	if !gotIpv4 && d.useUPnP {
		tools.Debug("get ipv4 address from upnp:")
		if ip, err := upnpGetPublicIp(); ip != nil {
			gotIpv4 = true
			ret = append(ret, ip.String())
			tools.Debug("  -> ipv4: %v", ip.String())
		} else if err == nil {
			tools.Debug("  -> not support")
		} else {
			tools.Debug("  -> error: %v", err)
		}
	}

	if !gotIpv4 && len(d.api) > 0 {
		tools.Debug("get ipv4 address from http (%v):", d.api)
		if ip, err := httpGetPublicIp(d.api); ip != nil {
			gotIpv4 = true
			ret = append(ret, ip.String())
			tools.Debug("  -> ipv4: %v", ip.String())
		} else if err == nil {
			tools.Debug("  -> no ip")
		} else {
			tools.Debug("  -> error: %v", err)
		}
	}

	tools.Debug("IP Address: %v", ret)

	d.lastGet = ret
}

func IsPublicIp(ip net.IP) bool {
	for _, n := range privateAddress {
		if n.Contains(ip) {
			return false
		}
	}
	return true
}

var privateAddress []*net.IPNet = nil

func init() {
	var _privateAddress = []string{
		// https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
		"192.88.99.0/24",
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.0.0/29",
		"192.0.0.8/32",
		"192.0.0.170/32",
		"192.0.0.171/32",
		"192.0.2.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32",
		"127.0.0.0/8",
		// https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml#note2
		"2001:10::/28",
		"::1/128",
		"::/128",
		"64:ff9b:1::/48",
		"100::/64",
		"2001:2::/48",
		"2001:db8::/32",
		"fe80::/10",
		"2001::/23",
		"fc00::/7",
		"2001::/32",
		"2002::/16",
	}

	privateAddress = make([]*net.IPNet, len(_privateAddress))
	for i, v := range _privateAddress {
		_, privateAddress[i], _ = net.ParseCIDR(v)
	}
}
