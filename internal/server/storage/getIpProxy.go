package storage

type iwrapGetIpOptions interface {
	GetIpHttpDisable() bool
	GetIpApi6() string
	GetIpApi4() string
}
type wrapGetIpOptions struct {
	iwrapGetIpOptions
}

func (w *wrapGetIpOptions) GetIpUpnpDisable() bool {
	return true
}
