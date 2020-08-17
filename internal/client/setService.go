package client

func (stat *ClientStateHolder) SetServices(services []string) {
	stat.statusData.lock()
	defer stat.statusData.unlock()

	stat.statusData.services = services
}
