package client

func (stat *ClientStateHolder) SetPublicPort(externalPort uint16) {
	stat.sharedStatus.lock()
	defer stat.sharedStatus.unlock()

	stat.sharedStatus.externalPort = uint32(externalPort)
}
