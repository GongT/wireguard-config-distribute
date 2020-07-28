package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/hostfile"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func main() {
	log.Println("program start.")
	opts := clientProgramOptions{}
	config.InitProgramArguments(&opts)

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
		opts.InternalIp = detect_ip.DetectLocalNetwork()
	}

	detect_ip.Detect(&opts.PublicIp, &opts.PublicIp6, !opts.GetIpHttpDsiable(), !opts.GetIpUpnpDsiable())

	tools.NormalizeServerString(&opts.Server)

	if opts.DebugMode {
		tools.Error("commandline arguments: %s", spew.Sdump(opts))
		tools.SetDebugMode(opts.DebugMode)
	}

	watcher := hostfile.StartWatch(opts.HostFile)
	c := client.NewClient(opts)
	c.ConfigureVPN(opts)
	c.Configure(opts)

	go func() {
		for content := range watcher.OnChange {
			c.SetServices(hostfile.ToArray(hostfile.ParseServices(content)))
		}
	}()

	c.StartCommunication()

	<-tools.WaitForCtrlC()

	watcher.StopWatch()
	c.Quit()

	fmt.Println("Bye, bye!")
}
