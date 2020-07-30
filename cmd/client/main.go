package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/elevate"
	"github.com/gongt/wireguard-config-distribute/internal/client/hostfile"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/detect_ip"
	"github.com/gongt/wireguard-config-distribute/internal/systemd"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func main() {
	spew.Config.Indent = "    "
	log.Println("program start.")
	opts := &clientProgramOptions{}
	config.InitProgramArguments(opts)

	if opts.DebugMode {
		tools.SetDebugMode(opts.DebugMode)
	}

	elevate.EnsureAdminPrivileges()

	if len(opts.Hostname) == 0 {
		opts.Hostname = os.Getenv("HOSTNAME")
	}
	if len(opts.Hostname) == 0 {
		opts.Hostname = os.Getenv("COMPUTERNAME")
	}
	if len(opts.Hostname) == 0 {
		tools.Die("HOSTNAME and COMPUTERNAME is empty, please set --hostname")
	}

	if len(opts.Title) == 0 {
		opts.Title = "Server at " + opts.Hostname
		tools.Error("title is not set, using hostname")
	}
	if opts.HostFile == "/etc/hosts" && runtime.GOOS == "windows" {
		opts.HostFile = "C:/Windows/System32/drivers/etc/hosts"
	}
	if len(opts.InternalIp) == 0 {
		ip, err := detect_ip.GetDefaultNetworkIp()
		if err != nil {
			tools.Die("Failed to find a valid local IP, please set --internal-ip")
		}
		opts.InternalIp = ip.String()
	}
	if len(opts.InterfaceName) == 0 {
		opts.InterfaceName = "wg_" + opts.JoinGroup
	}
	if len(opts.NetworkName) == 0 {
		mac, err := detect_ip.GetGatewayMac()
		if err != nil {
			tools.Error("Failed get gateway mac address: %s; using networking alone.", err.Error())
		}
		opts.NetworkName = "gw[" + strings.ToUpper(strings.ReplaceAll(mac, ":", "")) + "]"
	}

	detect_ip.Detect(&opts.PublicIp, &opts.PublicIp6, !opts.GetIpHttpDsiable(), !opts.GetIpUpnpDsiable())

	tools.NormalizeServerString(&opts.Server)

	if opts.DebugMode {
		tools.Error("commandline arguments: %s", spew.Sdump(opts))
	}

	watcher := hostfile.StartWatch(opts.HostFile)
	c := client.NewClient(opts)
	c.ConfigureVPN(opts)
	c.ConfigureInterface(opts)
	c.Configure(opts)

	go func() {
		for content := range watcher.OnChange {
			c.SetServices(hostfile.ToArray(hostfile.ParseServices(content)))
		}
	}()

	c.StartCommunication()

	systemd.ChangeToReady()
	<-tools.WaitForCtrlC()
	systemd.ChangeToQuit()

	watcher.StopWatch()
	c.Quit()

	fmt.Println("Bye, bye!")
}
