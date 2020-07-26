package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/config"
)

var opts = toolProgramOptions{}

func main() {
	parser := config.InitProgramArguments(&opts)
	spew.Dump(opts)

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
