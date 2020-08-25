package client

type oneTimeConfig struct {
	VpnGroupName        string
	Title               string
	Hostname            string
	LocalNetworkName    string
	ExternalEnabled     bool
	ExternalIp          []string
	ExternalPort        uint32
	InternalIp          string
	InternalPortDefault uint32
	InternalPort        uint32
	SelfMtu             uint16
}

func (cd *oneTimeConfig) configure(options configureOptions) {
	ip4 := options.GetPublicIp()
	if len(ip4) > 0 {
		cd.ExternalIp = append(cd.ExternalIp, ip4)
	}
	ip6 := options.GetPublicIp6()
	if len(ip6) > 0 {
		cd.ExternalIp = append(cd.ExternalIp, ip6)
	}

	if options.GetIpv6Only() && len(cd.ExternalIp) == 0 {
		cd.ExternalEnabled = false
	} else {
		cd.ExternalEnabled = true
	}

	cd.VpnGroupName = options.GetJoinGroup()
	cd.Title = options.GetTitle()
	cd.Hostname = options.GetHostname()
	cd.LocalNetworkName = options.GetNetworkName()
	cd.InternalIp = options.GetInternalIp()
	cd.SelfMtu = options.GetMTU()

	cd.ExternalPort = uint32(options.GetPublicPort())
	cd.InternalPortDefault = uint32(options.GetListenPort())
	cd.InternalPort = cd.InternalPortDefault
}
