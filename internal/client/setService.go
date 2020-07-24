package client

func (stat *clientStateHolder) SetServices(services []string) {
	fullServices := make([]string, len(services))
	for i, name := range services {
		fullServices[i] = name + "." + stat.configData.Hostname
	}

	stat.statusData.lock()
	stat.statusData.services = fullServices
	r := stat.isRunning
	stat.statusData.unlock()

	if r {
		stat.UploadInformation()
	}
}
