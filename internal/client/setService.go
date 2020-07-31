package client

func (stat *ClientStateHolder) SetServices(services []string) {
	stat.statusData.lock()

	stat.statusData.services = services
	r := stat.isRunning

	stat.statusData.unlock()

	if r {
		stat.uploadInformation()
	}
}
