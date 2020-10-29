package remoteControl

import (
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (tool *ToolObject) Dump() {
	ret, err := tool.server.DumpStatus()
	if err != nil {
		tools.Error("Server Error: %s", err)
		return
	}
	fmt.Println(ret.Text)
}
