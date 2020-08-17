package client

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/client/remoteControl"
	"github.com/gongt/wireguard-config-distribute/internal/constants"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (stat *ClientStateHolder) startNetwork() {
	// todo: try 5 times
	stat.server.Connect()
}

func (stat *ClientStateHolder) StartTool() *remoteControl.ToolObject {
	stat.startNetwork()

	return remoteControl.Create(stat.server)
}

func (stat *ClientStateHolder) StartCommunication() {
	stat.startNetwork()

	go func() {
		fmt.Println("Start communication...")
		for {
			if stat.isQuit {
				tools.Error("Event loop finished")
				return
			}

			stat.run()

			time.Sleep(5 * time.Second)
		}
	}()
}

func (stat *ClientStateHolder) run() {
	stat.isRunning = false

	tools.Error("Send handshake:")
	for {
		if stat.uploadInformation() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	tools.Error("Complete handshake")

	chanel, err := stat.server.Start(stat.sessionId.Serialize())
	if err != nil {
		tools.Error("grpc connected but start() failed, is server running? %s", err.Error())
		return
	}

	tmr := time.NewTicker(constants.KEEY_ALIVE_SECONDS)
	defer tmr.Stop()

	for {
		select {
		case <-tmr.C:
			if stat.isQuit {
				tools.Debug(" ~ quit")
				return
			}
			tools.Debug(" ~ send keep alive")
			result, err := stat.server.KeepAlive(stat.sessionId)
			if err != nil {
				tools.Error("grpc keep alive failed, is server (still) running? %s", err.Error())
				return
			}
			if !result.Success {
				tools.Error("server cleared, my state will reset.")
				return
			}
		case peers := <-chanel:
			if stat.isQuit {
				tools.Debug(" ~ quit")
				return
			} else if peers == nil {
				tools.Debug(" ~ server disconnected")
				return
			}
			tools.Debug(" ~ receive peers (%d peer, %d host)", len(peers.List), len(peers.Hosts))
			stat.vpn.UpdatePeers(peers.List)
			if stat.hostsHandler != nil {
				stat.hostsHandler(peers.Hosts)
			} else {
				tools.Error("hosts handler did not register?")
			}
		case <-stat.quitChan:
			tools.Debug(" ~ quit")
			return
		}
	}
}

type HandlerFunction = func(map[string]string)

func (stat *ClientStateHolder) HandleHosts(fn HandlerFunction) {
	stat.hostsHandler = fn
}
