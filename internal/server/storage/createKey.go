package storage

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"
)

const serverCertFileName = "ca.cert.pem"
const serverKeyFileName = "ca.key.pem"

func (storage *ServerStorage) createKey(ipList []net.IP, certFile, keyFile string) (err error) {
	fmt.Println("Siging certificate key...")
	if storage._cacheCaPri == nil || storage._cacheCa == nil {
		return errors.New("Invalid program state")
	}

	certPrivKey, err := readPKCS1(keyFile)
	if err != nil {
		certPrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}
		err = writePKCS1(keyFile, certPrivKey)
		if err != nil {
			fmt.Printf("  * Private key write failed: %s\n", err.Error())
			return
		}
		fmt.Printf("  * Private key has written to %s\n", keyFile)
	} else {
		fmt.Printf("  * Private key is read from %s\n", keyFile)
	}

	serverCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      storage._cacheCa.Subject,
		IPAddresses:  ipList,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 1, 4, 5, 1, 4},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, serverCert, storage._cacheCa, &certPrivKey.PublicKey, storage._cacheCaPri)
	if err != nil {
		return err
	}

	err = writeCert(certFile, certBytes)
	if err != nil {
		fmt.Printf("  * Public key write failed: %s\n", err.Error())
		return
	}
	fmt.Printf("  * Public key has written to %s\n", certFile)

	return
}
