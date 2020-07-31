package main

import (
	"errors"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (*clientProgramOptions) Sanitize() error {
	if !tools.GetSystemHostName(&opts.Hostname) {
		return errors.New("HOSTNAME and COMPUTERNAME is empty, please set --hostname")
	}

	if len(opts.Title) == 0 {
		opts.Title = "Server at " + opts.Hostname
		tools.Error("title is not set, using hostname")
	}
	if opts.HostFile == "/etc/hosts" && runtime.GOOS == "windows" {
		opts.HostFile = "C:/Windows/System32/drivers/etc/hosts"
	}
	if len(opts.InternalIp) == 0 {
		ip, err := detect_ip.GetDefaultNetworkIp()
		if err != nil {
			return errors.New("Failed to find a valid local IP, please set --internal-ip")
		}
		opts.InternalIp = ip.String()
	}
	if len(opts.InterfaceName) == 0 {
		opts.InterfaceName = "wg_" + opts.JoinGroup
	}
	if len(opts.NetworkName) == 0 {
		mac, err := detect_ip.GetGatewayMac()
		if err != nil {
			tools.Error("Failed get gateway mac address: %s; using networking alone.", err.Error())
		}
		opts.NetworkName = "gw[" + strings.ToUpper(strings.ReplaceAll(mac, ":", "")) + "]"
	}
	if len(opts.MachineID) == 0 {
		opts.MachineID, _ = machineid.ID()
	}

	detect_ip.Detect(&opts.PublicIp, &opts.PublicIp6, !opts.GetIpHttpDsiable(), !opts.GetIpUpnpDsiable())

	tools.NormalizeServerString(&opts.Server)

	return nil
}
