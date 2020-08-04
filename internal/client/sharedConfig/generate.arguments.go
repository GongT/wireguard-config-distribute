
package sharedConfig

type ReadOnlyConnectionOptions interface {
	GetServer() string
	GetPassword() string
	GetGrpcInsecure() bool
	GetGrpcHostname() string
	GetGrpcServerKey() string
}

func (self ConnectionOptions) GetServer() string {
	return self.Server
}

func (self ConnectionOptions) GetPassword() string {
	return self.Password
}

func (self ConnectionOptions) GetGrpcInsecure() bool {
	return self.GrpcInsecure
}

func (self ConnectionOptions) GetGrpcHostname() string {
	return self.GrpcHostname
}

func (self ConnectionOptions) GetGrpcServerKey() string {
	return self.GrpcServerKey
}

