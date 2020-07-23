package main

import (
	"fmt"
	"runtime"

	"github.com/gongt/wireguard-config-distribute/internal/client"
	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func main() {
	opts := clientProgramOptions{}
	config.ParseProgramArguments(&opts)

	if opts.HostFile == "/etc/hosts" && runtime.GOOS == "windows" {
		opts.HostFile = "C:/Windows/System32/drivers/etc/hosts"
	}
	opts.InternalIp

	creds, err := client.CreateClientTls(opts)
	if err != nil {
		tools.Die("Failed create TLS: %s", err.Error())
	}
	c := client.NewClient(opts, creds)

	c.StartNetwork()

	c.StartCommunication()

	tools.WaitForCtrlC()

	c.Quit()

	fmt.Println("Bye, bye!")
}
