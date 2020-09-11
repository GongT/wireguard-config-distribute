package remoteControl

import "fmt"

func (tool *ToolObject) Dump() {
	ret, err := tool.server.DumpStatus()
	if err != nil {
		panic(err)
	}
	fmt.Println(ret.Text)
}
