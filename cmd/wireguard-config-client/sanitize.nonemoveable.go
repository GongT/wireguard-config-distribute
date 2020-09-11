// +build !moveable

package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/upnp"
)

func (opts *clientProgramOptions) Sanitize() error {
	if err := opts.sanitizeBase(); err != nil {
		return err
	}

	if opts.GetNoPublicNetwork() {
		opts.PublicIp = []string{}
		opts.Gateway = false
		opts.IpUpnpDisable = true
		opts.IpApi4 = ""
		opts.IpApi6 = ""
		opts.NoAutoForwardUpnp = true
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
	if len(opts.NetworkName) == 0 {
		tools.Debug("detect gateway mac...")
		mac, err := detect_ip.GetGatewayMac()
		if err != nil {
			return fmt.Errorf("Failed get gateway mac address: %s", err.Error())
		}
		mac = strings.ReplaceAll(mac, ":", "")
		mac = strings.ReplaceAll(mac, "-", "")
		mac = strings.ToUpper(mac)
		opts.NetworkName = "gw[" + mac + "]"
		tools.Debug("  -> %s", mac)
	}

	if opts.PublicPort == 0 {
		opts.PublicPort = opts.ListenPort
	}
	if opts.VpnIpv6Only {
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

	return nil
}
