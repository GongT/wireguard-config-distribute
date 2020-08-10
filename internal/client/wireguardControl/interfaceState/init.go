package interfaceState

type InterfaceState interface {
	DeleteInterface() error
	CreateOrUpdateInterface(options InterfaceOptions) error
}

type interfaceState struct {
	ifname string

	address string
	network string
	mtu     int
	dns     []string
	table   string

	native *nativeState
}

// appendLineExtened("DNS = 1.1.1.1,8.8.8.8")
// appendLineExtened("Table = 12345")
// PreUp, PostUp, PreDown, PostDown
// SaveConfig?

func CreateInterface(ifname string) *interfaceState {
	ret := &interfaceState{
		ifname: ifname,
	}
	ret.init()
	return ret
}
