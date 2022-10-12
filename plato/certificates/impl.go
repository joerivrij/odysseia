package certificates

import (
	"crypto/rsa"
	"crypto/x509"
)

type CertClient interface {
	InitCa() error
	GenerateKeyAndCertSet(hosts []string, validityInDays int) ([]byte, []byte, error)
	PemEncodedCa() []byte
}

type CertificateGenerator struct {
	CaValidity    int
	Ca            *x509.Certificate
	CaPrivateKey  *rsa.PrivateKey
	CaPem         []byte
	CaPrivKeyPem  []byte
	Organizations []string
}

func NewCertGeneratorClient(organizations []string, caValidity int) (CertClient, error) {
	return &CertificateGenerator{
		CaValidity:    caValidity,
		Ca:            nil,
		CaPrivateKey:  nil,
		CaPem:         nil,
		CaPrivKeyPem:  nil,
		Organizations: organizations,
	}, nil
}
