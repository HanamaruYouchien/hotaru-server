package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ParsePrivateKey(keyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPEM)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ParsePublicKey(keyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(keyPEM)
	return x509.ParsePKCS1PublicKey(block.Bytes)
}
