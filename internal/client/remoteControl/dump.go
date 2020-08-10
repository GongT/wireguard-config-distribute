package remoteControl

func (tool *ToolObject) Dump() {
	ret, err := tool.server.DumpStatus()
	if err != nil {
		panic(err)
	}
	println(ret.Text)
}
