package remoteControl

import server "github.com/gongt/wireguard-config-distribute/internal/client/client.server"

type ToolObject struct {
	server server.ServerStatus
}

func Create(server server.ServerStatus) *ToolObject {
	ret := ToolObject{
		server: server,
	}
	return &ret
}
