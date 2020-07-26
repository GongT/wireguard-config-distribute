package client

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/client/remoteControl"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (stat *clientStateHolder) startNetwork() {
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

			stat.isRunning = false

			tools.Error("Send handshake:")
			for {
				if stat.UploadInformation() {
					break
				}
				time.Sleep(5 * time.Second)
			}

			// todo: boot vpn

			stat.tick()
		}
	}()
}

func (stat *clientStateHolder) tick() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := stat.server.KeepAlive(stat.SessionId)
			if err != nil || !result.Success {
				tools.Error("grpc keep alive failed, is server (still) running? %s", err.Error())
				return
			}
		case <-stat.quitChan:
			return
		}
	}
}
