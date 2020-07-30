package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

var opts = &toolProgramOptions{}

func main() {
	spew.Config.Indent = "    "
	parser := config.InitProgramArguments(opts)

	if opts.DebugMode {
		tools.Error("commandline arguments: %s", spew.Sdump(opts))
		tools.SetDebugMode(opts.DebugMode)
	}

	tools.NormalizeServerString(&opts.Server)

	c := client.NewClient(connectionOptions{
		server: opts.GetServer(),
	})

	tool := c.StartTool()

	if parser.Exists("download-ca") {
		tool.GetCA(opts.GetPassword(), opts.DownloadCA.GetOutput())
	} else {
		parser.DieUsage()
	}
}
