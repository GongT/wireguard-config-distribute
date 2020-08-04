package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type TLSOptions interface {
	GetGrpcInsecure() bool
	GetGrpcHostname() string
	GetGrpcServerKey() string
}

func createClientTls(opts TLSOptions) (credentials.TransportCredentials, error) {
	cfg := tls.Config{}

	if opts.GetGrpcInsecure() {
		cfg.InsecureSkipVerify = true
	}
	if len(opts.GetGrpcServerKey()) > 0 {
		b, err := ioutil.ReadFile(opts.GetGrpcServerKey())
		if err != nil {
			return nil, err
		}
		cp, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		if !cp.AppendCertsFromPEM(b) {
			return nil, errors.New("credentials: failed to append certificates")
		}
		cfg.RootCAs = cp
	}
	if len(opts.GetGrpcHostname()) > 0 {
		cfg.ServerName = opts.GetGrpcHostname()
	}

	return credentials.NewTLS(&cfg), nil
}
