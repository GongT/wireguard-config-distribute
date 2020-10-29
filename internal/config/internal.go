package config

type commonOptions struct {
	DebugMode   bool `long:"debug" short:"D" description:"enable debug mode" env:"WIREGUARD_CONFIG_DEVELOPMENT"`
	ShowVersion bool `long:"version" short:"V" description:"show program version and exit"`
}

type ReadOnlyCommonOptions interface {
	GetDebugMode() bool
	GetShowVersion() bool
}

func (self commonOptions) GetDebugMode() bool {
	return self.DebugMode
}

func (self commonOptions) GetShowVersion() bool {
	return self.ShowVersion
}
