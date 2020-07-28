package client

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/client/remoteControl"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (stat *clientStateHolder) startNetwork() {
	// todo: try 5 times
	stat.server.Connect()
}

func (stat *clientStateHolder) StartTool() *remoteControl.ToolObject {
	stat.startNetwork()

	return remoteControl.Create(stat.server)
}

func (stat *clientStateHolder) StartCommunication() {
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

func (stat *clientStateHolder) run() {
	stat.isRunning = false

	tools.Error("Send handshake:")
	for {
		if stat.uploadInformation() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	tools.Error("Complete handshake")

	channel, err := stat.server.Start(stat.MachineId)
	if err != nil {
		tools.Error("grpc connected but start() failed, is server running? %s", err.Error())
		return
	}

	tmr := time.NewTicker(20 * time.Second)
	defer tmr.Stop()

	for {
		select {
		case <-tmr.C:
			if stat.isQuit {
				tools.Debug(" ~ quit")
				return
			}
			tools.Debug(" ~ send keep alive")
			result, err := stat.server.KeepAlive(stat.MachineId)
			if err != nil {
				tools.Error("grpc keep alive failed, is server (still) running? %s", err.Error())
				return
			}
			if !result.Success {
				tools.Error("server cleared, my state will reset.")
				return
			}
		case peers := <-channel:
			if stat.isQuit {
				tools.Debug(" ~ quit")
				return
			} else if peers == nil {
				tools.Debug(" ~ server disconnected")
				return
			}
			tools.Debug(" ~ receive peers (%d peer)", len(peers.List))
			for _, peer := range peers.List {
				tools.Debug("  * <%s> %s -> %s", peer.MachineId, peer.Hostname, peer.GetPeer().GetAddress())
			}
		case <-stat.quitChan:
			tools.Debug(" ~ quit")
			return
		}
	}
}
