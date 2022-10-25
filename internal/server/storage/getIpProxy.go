package storage

type iwrapGetIpOptions interface {
	GetIpHttpDisable() bool
	GetIpApi() string
}
type wrapGetIpOptions struct {
	iwrapGetIpOptions
}

func (w *wrapGetIpOptions) GetIpUpnpDisable() bool {
	return true
}
