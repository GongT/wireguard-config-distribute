package client

func (stat *ClientStateHolder) SetServices(services []string) {
	stat.sharedStatus.lock()
	defer stat.sharedStatus.unlock()

	stat.sharedStatus.services = services
}
