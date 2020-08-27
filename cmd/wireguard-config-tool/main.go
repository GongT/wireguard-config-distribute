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

	if config.CommandActive("download-ca") {
		opts.ConnectionOptions.GrpcInsecure = true
		c().GetCA(opts.DownloadCA.GetOutput())
	} else if config.CommandActive("netgroup") {
		optsNg := opts.GetNetworkGroup()
		if config.CommandActive("create") {
			c().CreateNetworkGroup(optsNg.GetCreate())
		} else if config.CommandActive("create") {
			c().DeleteNetworkGroup(optsNg.GetDelete())
		} else {
			config.DieUsage()
		}
	} else if config.CommandActive("dump") {
		c().Dump()
	} else {
		config.DieUsage()
	}
}

func c() *remoteControl.ToolObject {
	return client.NewClient(opts.GetConnectionOptions()).StartTool()
}
