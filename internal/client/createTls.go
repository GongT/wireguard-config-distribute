package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type tlsClientOptions interface {
	GetGrpcInsecure() bool
	GetGrpcHostname() string
	GetGrpcServerKey() string
}

func CreateClientTls(opts tlsClientOptions) (credentials.TransportCredentials, error) {
	cfg := tls.Config{}

	if opts.GetGrpcInsecure() {
		cfg.InsecureSkipVerify = true
	}
	if len(opts.GetGrpcServerKey()) > 0 {
		b, err := ioutil.ReadFile(opts.GetGrpcServerKey())
		if err != nil {
			return nil, err
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("credentials: failed to append certificates")
		}
		cfg.RootCAs = cp
	}
	if len(opts.GetGrpcHostname()) > 0 {
		cfg.ServerName = opts.GetGrpcHostname()
	}

	return credentials.NewTLS(&cfg), nil
}
