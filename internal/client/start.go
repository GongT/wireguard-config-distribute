package client

import (
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (stat *clientStateHolder) StartCommunication() {
	stat.server.Connect()

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		fmt.Println("start communication...")
		for {
			select {
			case <-ticker.C:
				_, err := stat.server.KeepAlive()
				if err != nil {
					tools.Error("grpc keep alive failed, is server running? %s", err.Error())
				}
			case <-stat.quitChan:
				ticker.Stop()
				fmt.Println("stop communication.")
				return
			}
		}
	}()

	for {
		if stat.UploadInformation() {
			break
		}
	}
}
