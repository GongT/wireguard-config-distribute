package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/hostfile"
	"github.com/gongt/wireguard-config-distribute/internal/client/service"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/systemd"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/upnp"
	"github.com/judwhite/go-svc/svc"
)

type program struct {
	client      *client.ClientStateHolder
	watcher     *hostfile.Watcher
	portForward *upnp.UPnPPortForwarder
}

var opts = &clientProgramOptions{}
var prog = &program{}

func main() {
	spew.Config.Indent = "    "
	dispose := tools.WaitExit(func(code int) {
		tools.Error("program dying!")
		svc.Service.Stop(prog)
		tools.Error("service stop complete.")
	})

	if err := svc.Run(prog); err != nil {
		tools.Error("Failed run service: %s", err.Error())
		tools.HasError(100)
	}

	dispose()
	tools.ExitMain()
}

func (p *program) Init(env svc.Environment) error {
	err := config.InitProgramArguments(opts)
	p.client = client.NewClient(opts.GetConnectionOptions())

	if err != nil {
		return err
	}

	if !env.IsWindowsService() {
		service.EnsureAdminPrivileges(opts)
	}

	return nil
}

func (p *program) Start() error {
	log.Println("service start.")

	p.watcher = hostfile.StartWatch(opts.HostFile)
	p.client.Configure(opts)

	go func() {
		for range p.watcher.OnChange {
		}
	}()
	// go func() {
	// 	for content := range p.watcher.OnChange {
	// 		p.client.SetServices(hostfile.ToArray(hostfile.ParseServices(content)))
	// 	}
	// }()

	if opts.GetNoAutoForwardUpnp() {
		tools.Debug("[UPnP] disabled")
	}else{
		tools.Debug("[UPnP] enabled")
		portForward, err := upnp.NewAutoForward(opts)
		if err != nil {
			tools.Die("Failed init upnp, you may disable it if you do not use: %v")
		}

		p.portForward = portForward
		port, err := p.portForward.Start()
		if err != nil {
			tools.Error("[UPNP] Forward first tick error: %v", err)
		} else {
			p.client.SetPublicPort(port)
			go func() {
				for port := range p.portForward.OnChange {
					p.client.SetPublicPort(port)
				}
			}()
		}
	}

	p.client.HandleHosts(func(hosts map[string]string) {
		p.watcher.WriteBlock(hosts)
	})

	p.client.StartCommunication()

	return nil
}

func (p *program) Stop() error {
	fmt.Println("Service is quitting!")

	systemd.ChangeToQuit()

	if p.watcher != nil {
		p.watcher.StopWatch()
	}
	if p.portForward != nil {
		p.portForward.Close()
	}

	p.client.Quit()

	fmt.Println("Bye, bye!")

	return nil
}
