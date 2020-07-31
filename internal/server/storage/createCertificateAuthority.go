package storage

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

const caCertFileName = "ca.cert.pem"
const caKeyFileName = "ca.key.pem"

func (storage *ServerStorage) loadOrCreateCA(serverName string) (ca *x509.Certificate, pri *rsa.PrivateKey, err error) {
	if storage._cacheCaPri != nil {
		return storage._cacheCa, storage._cacheCaPri, nil
	}

	ca, err = readCert(storage.Path(caCertFileName))
	if err != nil || ca.Subject.CommonName != serverName {
		fmt.Println("Creating self-signed TLS certificate authority (CA)")
		err = storage.createSelfCA(serverName)
		if err != nil {
			return
		}
		ca, err = readCert(storage.Path(caCertFileName))
		if err != nil {
			return
		}
	}

	pri, err = readPKCS1(storage.Path(caKeyFileName))
	if err != nil {
		return
	}

	fmt.Println("Self-signed TLS certificate authority (CA) loaded")
	storage._cacheCa = ca
	storage._cacheCaPri = pri

	return
}

func (storage *ServerStorage) createSelfCA(serverName string) (err error) {
	caFile := storage.CaCertFilePath()
	caPriFile := storage.caKeyFilePath()

	caCert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{serverName},
			CommonName:   serverName,
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
		fmt.Printf("  * Public key write failed: %s\n", err.Error())
		return
	}
	fmt.Printf("  * Public CA key has written to %s\n", caFile)

	err = writePKCS1(caPriFile, caPrivKey)
	if err != nil {
		fmt.Printf("  * Private key write failed: %s\n", err.Error())
		return
	}
	fmt.Printf("  * Private CA key has written to %s\n", caPriFile)

	return
}

func (storage *ServerStorage) CaCertFilePath() string {
	return storage.Path(caCertFileName)
}

func (storage *ServerStorage) caKeyFilePath() string {
	return storage.Path(caKeyFileName)
}

func (storage *ServerStorage) GetCaCertFileContent() []byte {
	f := storage.CaCertFilePath()
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		tools.Error("Failed read file %s: %s", f, err.Error())
		return nil
	}

	return bs
}
