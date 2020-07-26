package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type TLSOptions struct {
	Insecure  bool
	Hostname  string
	ServerKey string
}

func createClientTls(opts TLSOptions) (credentials.TransportCredentials, error) {
	cfg := tls.Config{}

	if opts.Insecure {
		cfg.InsecureSkipVerify = true
	}
	if len(opts.ServerKey) > 0 {
		b, err := ioutil.ReadFile(opts.ServerKey)
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
	if len(opts.Hostname) > 0 {
		cfg.ServerName = opts.Hostname
	}

	return credentials.NewTLS(&cfg), nil
}
