// +build moveable

package main

func (opts *clientProgramOptions) SanitizeBase() error {
	return opts.sanitizeBase()
}

func (self clientProgramOptions) GetIpv4Only() bool {
	return false
}

func (self clientProgramOptions) GetIpv6Only() bool {
	return false
}

func (self clientProgramOptions) GetPublicIp() string {
	return ""
}

func (self clientProgramOptions) GetPublicIp6() string {
	return ""
}

func (self clientProgramOptions) GetPublicPort() uint16 {
	return 0
}

func (self clientProgramOptions) GetIpApi6() string {
	return ""
}

func (self clientProgramOptions) GetIpApi4() string {
	return ""
}

func (self clientProgramOptions) GetIpUpnpDisable() bool {
	return true
}

func (self clientProgramOptions) GetIpHttpDisable() bool {
	return true
}

func (self clientProgramOptions) GetNoAutoForwardUpnp() bool {
	return true
}

func (self clientProgramOptions) GetInternalIp() string {
	return ""
}

func (self clientProgramOptions) NoPublicNetwork() bool {
	return true
}
