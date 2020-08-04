package main

import "github.com/gongt/wireguard-config-distribute/internal/tools"

func (opts *toolProgramOptions) Sanitize() error {
	tools.NormalizeServerString(&opts.ConnectionOptions.Server)
	return nil
}
