package main

import (
	"fmt"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/hostfile"
	"github.com/gongt/wireguard-config-distribute/internal/client/network_detect"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/upnp"
)

func main() {
	opts := clientProgramOptions{}
	config.ParseProgramArguments(&opts)

	if len(opts.Hostname) == 0 {
		tools.Die("HOSTNAME is empty, please set --hostname")
	}

	if len(opts.Title) == 0 {
		opts.Title = "Server at " + opts.Hostname
		tools.Error("title is not set, using hostname")
	}
	if opts.HostFile == "/etc/hosts" && runtime.GOOS == "windows" {
		opts.HostFile = "C:/Windows/System32/drivers/etc/hosts"
	}
	if len(opts.InternalIp) == 0 {
		opts.InternalIp = network_detect.DetectLocalNetwork()
	}
	if len(opts.PublicIp) == 0 && !opts.Ipv6Only {
		if !opts.IpUpnpDsiable {
			tools.Error("Finding public ipv4 from router...")
			opts.PublicIp, _ = upnp.GetPublicIp()
		}
		if len(opts.PublicIp) == 0 && !opts.IpHttpDsiable {
			tools.Error("Fetching IPv4 address...")
			opts.PublicIp, _ = detect_ip.GetPublicIp()
		}
	}
	if len(opts.PublicIp6) == 0 {
		if len(opts.PublicIp) == 0 && !opts.IpHttpDsiable {
			tools.Error("Fetching IPv6 address...")
			opts.PublicIp6, _ = detect_ip.GetPublicIp6()
		}
	}

	tools.Debug("input config: %s", spew.Sdump(opts))

	if opts.DebugMode {
		tools.SetDebugMode(opts.DebugMode)
	}

	watcher := hostfile.StartWatch(opts.HostFile)
	c := client.NewClient(opts)
	c.Configure(opts)

	go func() {
		for content := range watcher.OnChange {
			c.SetServices(hostfile.ToArray(hostfile.ParseServices(content)))
		}
	}()

	c.StartCommunication()

	<-tools.WaitForCtrlC()

	c.Quit()
	watcher.StopWatch()

	fmt.Println("Bye, bye!")
}
