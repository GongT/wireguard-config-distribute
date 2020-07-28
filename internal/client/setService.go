package client

func (stat *clientStateHolder) SetServices(services []string) {
	stat.statusData.lock()

	stat.statusData.services = services
	r := stat.isRunning

	stat.statusData.unlock()

	if r {
		stat.uploadInformation()
	}
}
