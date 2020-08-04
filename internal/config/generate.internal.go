
package config

type ReadOnlyInternalOptions interface {
	GetStandardOutputPath() string
	GetIsElevated() bool
}

func (self internalOptions) GetStandardOutputPath() string {
	return self.StandardOutputPath
}

func (self internalOptions) GetIsElevated() bool {
	return self.IsElevated
}

type ReadOnlyCommonOptions interface {
	GetDebugMode() bool
	GetLogFilePath() string
}

func (self commonOptions) GetDebugMode() bool {
	return self.DebugMode
}

func (self commonOptions) GetLogFilePath() string {
	return self.LogFilePath
}

