package client

import "fmt"

type infoOptions interface {
	GetJoinGroup() string
	GetInternalIp() string
	GetListenPort() uint16
	GetPublicPort() uint16
	GetMTU() uint16
	GetTitle() string
	GetHostname() string
	GetNetworkName() string
}

type infoConfig struct {
	VpnGroupName        string
	Title               string
	Hostname            string
	LocalNetworkName    string
	ExternalPort        uint32
	InternalIp          string
	InternalPortDefault uint32
	InternalPort        uint32
	SelfMtu             uint16
}

func (cd *infoConfig) createInterfaceComment() string {
	return fmt.Sprintf("%s (%s) [AT] %s", cd.Hostname, cd.Title, cd.LocalNetworkName)
}

func (cd *infoConfig) configure(options infoOptions) {
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
