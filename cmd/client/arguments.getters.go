package main

type readonlyClientProgramOptions interface {
	GetListenPort() uint16
	GetMTU() uint16
	GetServer() string
	GetNetworkName() string
	GetJoinGroup() string
	GetPerferIp() string
	GetTitle() string
	GetHostname() string
	GetHostFile() string
	GetIpv6Only() bool
	GetPublicIp() string
	GetIpServerDsiable() bool
	GetIpUpnpDsiable() bool
	GetIpHttpDsiable() bool
	GetInternalIp() string
	GetGrpcInsecure() bool
	GetGrpcHostname() string
	GetGrpcServerKey() string
	GetDebugMode() bool
}

func (s clientProgramOptions) GetListenPort() uint16 {
	return s.ListenPort;
}

func (s clientProgramOptions) GetMTU() uint16 {
	return s.MTU;
}

func (s clientProgramOptions) GetServer() string {
	return s.Server;
}

func (s clientProgramOptions) GetNetworkName() string {
	return s.NetworkName;
}

func (s clientProgramOptions) GetJoinGroup() string {
	return s.JoinGroup;
}

func (s clientProgramOptions) GetPerferIp() string {
	return s.PerferIp;
}

func (s clientProgramOptions) GetTitle() string {
	return s.Title;
}

func (s clientProgramOptions) GetHostname() string {
	return s.Hostname;
}

func (s clientProgramOptions) GetHostFile() string {
	return s.HostFile;
}

func (s clientProgramOptions) GetIpv6Only() bool {
	return s.Ipv6Only;
}

func (s clientProgramOptions) GetPublicIp() string {
	return s.PublicIp;
}

func (s clientProgramOptions) GetIpServerDsiable() bool {
	return s.IpServerDsiable;
}

func (s clientProgramOptions) GetIpUpnpDsiable() bool {
	return s.IpUpnpDsiable;
}

func (s clientProgramOptions) GetIpHttpDsiable() bool {
	return s.IpHttpDsiable;
}

func (s clientProgramOptions) GetInternalIp() string {
	return s.InternalIp;
}

func (s clientProgramOptions) GetGrpcInsecure() bool {
	return s.GrpcInsecure;
}

func (s clientProgramOptions) GetGrpcHostname() string {
	return s.GrpcHostname;
}

func (s clientProgramOptions) GetGrpcServerKey() string {
	return s.GrpcServerKey;
}

func (s clientProgramOptions) GetDebugMode() bool {
	return s.DebugMode;
}


