//go:generate go-generate-struct-interface

package main

import "github.com/gongt/wireguard-config-distribute/internal/client/sharedConfig"

type downloadCaAction struct {
	Output string `short:"O" long:"output" description:"output file path" default:"./server.cert"`
}
type netgroupAction struct {
	Create netgroupCreateAction `command:"create" description:"create network group"`
	Delete netgroupDeleteAction `command:"delete" description:"delete network group"`
}
type netgroupCreateAction struct {
	Name   string `short:"n" long:"name" description:"name of the group, name of interfaces, only allow [0-9a-z_]" required:"true"`
	Prefix string `short:"p" long:"prefix" description:"ip prefix, eg: 192.168.100" required:"true"`
	Title  string `short:"t" long:"title" description:"friendly name of this group"`
}
type netgroupDeleteAction struct {
	Name string `short:"n" long:"name" description:"name of the group" required:"true"`
}

type toolProgramOptions struct {
	ConnectionOptions sharedConfig.ConnectionOptions `group:"Connection Options"`

	DownloadCA   downloadCaAction `command:"download-ca" alias:"auth" description:"Download server's self-signed CA cert file"`
	NetworkGroup netgroupAction   `command:"netgroup" description:"configure VPN network groups"`
}
