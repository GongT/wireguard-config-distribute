package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/client/hostfile"
	"github.com/gongt/wireguard-config-distribute/internal/client/service"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/systemd"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/judwhite/go-svc/svc"
)

type program struct {
	client  *client.ClientStateHolder
	watcher *hostfile.Watcher
}

var opts = &clientProgramOptions{}

func main() {
	prg := &program{}

	go func() {
		v := <-tools.WaitDie()
		tools.Error("program dying! (%s)", v)
		svc.Service.Stop(prg)
		tools.Error("service stop complete.")
	}()

	if err := svc.Run(prg); err != nil {
		tools.Error("Failed run service: %s", err.Error())
		os.Exit(1)
	}
}

func (p *program) Init(env svc.Environment) error {
	log.Println("program start.")

	spew.Config.Indent = "    "
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
	p.watcher = hostfile.StartWatch(opts.HostFile)
	p.client.ConfigureVPN(opts)
	p.client.Configure(opts)

	go func() {
		for content := range p.watcher.OnChange {
			p.client.SetServices(hostfile.ToArray(hostfile.ParseServices(content)))
		}
	}()

	p.client.StartCommunication()

	systemd.ChangeToReady()

	return nil
}

func (p *program) Stop() error {
	fmt.Println("Service is quitting!")

	systemd.ChangeToQuit()

	p.watcher.StopWatch()
	p.client.Quit()

	config.Cleanup()

	fmt.Println("Bye, bye!")

	os.Exit(0)

	return nil
}
