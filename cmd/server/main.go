package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/server"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

var listenPath *string = flag.String("unix", "", "listen unix socket")
var listenPort *uint = flag.Uint("port", 0, "listen unix socket")

func main() {
	fmt.Println("Hello, Server!")

	flag.Parse()

	s := server.NewServer()

	if len(*listenPath) > 0 || *listenPort > 0 {
		if len(*listenPath) > 0 {
			s.ListenSocket(server.ListenUnix(*listenPath))
		}
		if *listenPort > 0 {
			s.ListenSocket(server.ListenTCP(uint16(*listenPort)))
		}
	} else {
		listen := config.GetConfig(config.CONFIG_SERVER_LISTEN, config.CONFIG_SERVER_LISTEN_DEFAULT)
		v, err := strconv.ParseUint(listen, 10, 16)
		if err == nil {
			s.ListenSocket(server.ListenTCP(uint16(v)))
		} else {
			s.ListenSocket(server.ListenUnix(listen))
		}
	}

	tools.WaitForCtrlC()
	fmt.Println("Bye, bye!")
}
