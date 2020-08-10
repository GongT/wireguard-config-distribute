package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/remoteControl"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

var opts = &toolProgramOptions{}

func main() {
	spew.Config.Indent = "    "
	err := config.InitProgramArguments(opts)

	if err != nil {
		tools.Die("invalid commandline arguments: %s", err.Error())
	}

	if config.Exists("download-ca") {
		opts.ConnectionOptions.GrpcInsecure = true
		c().GetCA(opts.DownloadCA.GetOutput())
	} else if config.Exists("netgroup") {
		optsNg := opts.GetNetworkGroup()
		if config.Exists("create") {
			c().CreateNetworkGroup(optsNg.GetCreate())
		} else if config.Exists("create") {
			c().DeleteNetworkGroup(optsNg.GetDelete())
		} else {
			config.DieUsage()
		}
	} else if config.Exists("dump") {
		c().Dump()
	} else {
		config.DieUsage()
	}
}

func c() *remoteControl.ToolObject {
	return client.NewClient(opts.GetConnectionOptions()).StartTool()
}
