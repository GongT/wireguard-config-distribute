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
var logger *os.File

func main() {
	prg := &program{}
	if err := svc.Run(prg); err != nil {
		tools.Error("Failed run service: %s", err.Error())
		os.Exit(1)
	}
}

func init() {
	log.Println("program init.")
	config.UpdateDebug = func(debug bool) {
		tools.SetDebugMode(debug)

		if f := opts.GetLogFilePath(); len(f) > 0 {
			log.Println("log will dup to file.")
			logger = service.SetLogOutput(f)
			log.Println("log start.")
		}
	}
}

func (p *program) Init(env svc.Environment) error {
	log.Println("program start.")

	spew.Config.Indent = "    "
	_, err := config.InitProgramArguments(opts)
	p.client = client.NewClient(opts)

	if opts.GetDebugMode() {
		tools.Error("commandline arguments: %s", spew.Sdump(opts))
	}

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
	p.client.ConfigureInterface(opts)
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

	fmt.Println("Bye, bye!")

	if logger != nil {
		if err := logger.Sync(); err != nil {
			tools.Error("file.Sync() fail: %s", err.Error())
		}
		if err := logger.Close(); err != nil {
			tools.Error("file.Close() fail: %s", err.Error())
		}
	}

	return nil
}
