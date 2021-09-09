package remoteControl

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type createOptions interface {
	GetName() string
	GetPrefix() string
	GetTitle() string
}

func (tool *ToolObject) CreateNetworkGroup(opts createOptions) {
	err := tool.server.NewGroup(&protocol.NewGroupRequest{
		Name:     opts.GetName(),
		IpPrefix: opts.GetPrefix(),
		Title:    opts.GetTitle(),
	})
	if err == nil {
		tools.Error("successfully created VPN network.")
	} else {
		tools.Error("failed create VPN network: %v", err)
	}
}

type deleteOptions interface {
	GetName() string
}

func (tool *ToolObject) DeleteNetworkGroup(opts deleteOptions) {

}
