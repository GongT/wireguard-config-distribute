package storage

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

const caCertFileName = "ca.cert.pem"
const caKeyFileName = "ca.key.pem"

func (storage *ServerStorage) getCA(serverName string) (*x509.Certificate, []byte, error) {

}

func (storage *ServerStorage) createCA(serverName string) (err error) {
	caFile := storage.Path(caCertFileName)
	caPriFile := storage.Path(caKeyFileName)

	caCert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{serverName},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return
	}

	err = writeCert(caFile, caBytes)
	if err != nil {
		return
	}

	err = writePKCS1(caPriFile, caPrivKey)
	if err != nil {
		return
	}

	return
}
