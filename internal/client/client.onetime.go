package client

type oneTimeConfig struct {
	GroupName       string
	Title           string
	Hostname        string
	NetworkId       string
	ExternalEnabled bool
	ExternalIp      []string
	ExternalPort    uint32
	InternalIp      string
	InternalPort    uint32
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

	cd.GroupName = options.GetJoinGroup()
	cd.Title = options.GetTitle()
	cd.Hostname = options.GetHostname()
	cd.NetworkId = options.GetNetworkName()
	cd.ExternalPort = uint32(options.GetListenPort())
	cd.InternalIp = options.GetInternalIp()
	cd.InternalPort = uint32(options.GetListenPort())
}
