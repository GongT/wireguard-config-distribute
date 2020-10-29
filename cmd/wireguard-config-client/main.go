package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/autoUpdate"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/hostfile"
	"github.com/gongt/wireguard-config-distribute/internal/client/service"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/upnp"
)

var opts = &clientProgramOptions{}

func main() {
	spew.Config.Indent = "    "
	dispose := tools.WaitExit(func(code int) {
		tools.Error("program dying!")
		tools.Error("service stop complete.")
	})

	err := config.InitProgramArguments(opts)
	if err != nil {
		tools.HasFatalError()
		tools.Die("failed parse arguments: %s", err)
	}

	service.EnsureAdminPrivileges(opts)

	clientInstance := client.NewClient(opts.GetConnectionOptions())
	clientInstance.Configure(opts)

	go autoUpdate.StartAutoUpdate()

	log.Println("service start.")

	watcher := hostfile.StartWatch(opts.HostFile)
	go func() {
		//todo: this not work...
		/*
			for content := range watcher.OnChange {
				clientInstance.SetServices(hostfile.ToArray(hostfile.ParseServices(content)))
			}
		*/
	}()

	if opts.GetNoAutoForwardUpnp() {
		tools.Debug("[UPnP] disabled")
	} else {
		tools.Debug("[UPnP] enabled")
		portForward, err := upnp.NewAutoForward(opts)
		if err != nil {
			tools.Die("Failed init upnp, you may disable it if you do not use: %v")
		}

		port, err := portForward.Start()
		if err != nil {
			tools.Error("[UPNP] Forward first tick error: %v", err)
		} else {
			clientInstance.SetPublicPort(port)
			go func() {
				for port := range portForward.OnChange {
					clientInstance.SetPublicPort(port)
				}
			}()
		}
	}

	clientInstance.HandleHosts(func(hosts map[string]string) {
		watcher.WriteBlock(opts.GetJoinGroup(), hosts)
	})
	tools.WaitExit(func(int) {
		watcher.WriteBlock(opts.GetJoinGroup(), nil)
	})

	clientInstance.StartCommunication()

	<-tools.WaitForCtrlC()

	fmt.Println("Bye, bye.")
	dispose()
	tools.HasNoError()
	tools.ExitMain()

	tools.Die("this will never run")
}
