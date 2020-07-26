package storage

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func readPKCS1(source string) (key *rsa.PrivateKey, err error) {
	bytes, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	block, _ := pem.Decode(bytes)
	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("expect file " + source + " contains a rsa private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func writePKCS1(target string, content *rsa.PrivateKey) error {
	file, err := os.OpenFile(target, os.O_WRONLY+os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	err = pem.Encode(file, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(content),
	})
	if err != nil {
		return err
	}

	return nil
}

func readCert(source string) (*x509.Certificate, error) {
	bytes, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block.Type != "CERTIFICATE" {
		return nil, errors.New("expect file " + source + " contains a certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

func writeCert(target string, content []byte) error {
	err := os.MkdirAll(filepath.Dir(target), os.FileMode(0755))
	if err != nil {
		return err
	}

	file, err := os.OpenFile(target, os.O_WRONLY+os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}

	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: content,
	})
	if err != nil {
		return err
	}

	return nil
}
