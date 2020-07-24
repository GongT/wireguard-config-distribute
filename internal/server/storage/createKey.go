package storage

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

const serverCertFileName = "ca.cert.pem"
const serverKeyFileName = "ca.key.pem"

func (storage *ServerStorage) createKey(serverName, ipList []net.IP) (err error) {
	certFile := storage.Path(serverCertFileName)
	keyFile := storage.Path(serverKeyFileName)

	serverCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{serverName},
			Country:      []string{"US"},
		},
		IPAddresses:  ipList,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, serverCert, caData, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	err = writeCert(certFile, certBytes)
	if err != nil {
		return
	}

	err = writePKCS1(keyFile, certPrivKey)
	if err != nil {
		return
	}
}
