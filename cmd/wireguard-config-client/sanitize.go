package main

import (
	"errors"
	"fmt"

	"github.com/denisbrodbeck/machineid"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/phayes/freeport"
)

func (opts *clientProgramOptions) sanitizeBase() error {
	if !tools.GetSystemHostName(&opts.Hostname) {
		return errors.New("HOSTNAME and COMPUTERNAME is empty, please set --hostname")
	}

	if len(opts.Title) == 0 {
		opts.Title = "Server at " + opts.Hostname
		tools.Error("title is not set, using hostname")
	}
	if len(opts.InterfaceName) == 0 {
		opts.InterfaceName = "wg_" + opts.JoinGroup
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

	tools.NormalizeServerString(&opts.ConnectionOptions.Server)

	return nil
}
