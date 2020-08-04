package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/upnp"
	"github.com/phayes/freeport"
)

func (opts *clientProgramOptions) Sanitize() error {
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
		tools.Debug("detect default network source ip...")
		ip, err := detect_ip.GetDefaultNetworkIp()
		if err != nil {
			return errors.New("Failed to find a valid local IP, please set --internal-ip")
		}
		opts.InternalIp = ip.String()
		tools.Debug("  -> %s", opts.InternalIp)
	}
	if len(opts.InterfaceName) == 0 {
		opts.InterfaceName = "wg_" + opts.JoinGroup
	}
	if len(opts.NetworkName) == 0 {
		tools.Debug("detect gateway mac...")
		mac, err := detect_ip.GetGatewayMac()
		if err != nil {
			return fmt.Errorf("Failed get gateway mac address: %s", err.Error())
		}
		opts.NetworkName = "gw[" + strings.ToUpper(strings.ReplaceAll(mac, ":", "")) + "]"
		tools.Debug("  -> %s", mac)
	}
	if len(opts.MachineID) == 0 {
		tools.Debug("get machine guid...")
		opts.MachineID, _ = machineid.ID()
		tools.Debug("  -> %s", opts.MachineID)
	}

	if opts.ListenPort == 0 {
		tools.Debug("get free port...")
		port, err := freeport.GetFreePort()
		if err != nil {
			return fmt.Errorf("Failed find free port: %s", err.Error())
		}
		tools.Debug("  -> %d", port)
		opts.ListenPort = uint16(port)
	}
	if opts.PublicPort == 0 {
		opts.PublicPort = opts.ListenPort
	}
	if opts.Ipv6Only {
		opts.NoAutoForwardUpnp = true
	}
	if !opts.NoAutoForwardUpnp {
		tools.Debug("forward port with UPnP...")
		p, err := upnp.TryAddPortMapping(int(opts.ListenPort), int(opts.PublicPort))
		if err != nil {
			return fmt.Errorf("Failed forward port with UPnP: %s", err.Error())
		}
		tools.Debug("  -> %d", p)
		opts.PublicPort = p
	}

	detect_ip.Detect(&opts.PublicIp, &opts.PublicIp6, !opts.GetIpHttpDsiable(), !opts.GetIpUpnpDsiable())

	if len(opts.PublicIp) == 0 && !opts.Ipv6Only {
		return errors.New("Failed find an ipv4 address, and --ipv6only not set")
	}

	tools.NormalizeServerString(&opts.ConnectionOptions.Server)

	return nil
}
