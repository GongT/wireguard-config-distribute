package main

type connectionOptions struct {
	server string
}

func (c connectionOptions) GetServer() string { return c.server }
func (c connectionOptions) GetGrpcInsecure() bool {
	return true
}
func (c connectionOptions) GetGrpcHostname() string  { return "" }
func (c connectionOptions) GetGrpcServerKey() string { return "" }
