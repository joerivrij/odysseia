package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

func (c *CertificateGenerator) InitCa() error {
	ca, privateKey, err := c.generateCa()
	if err != nil {
		return err
	}

	c.CaPrivateKey = privateKey
	c.Ca = ca

	return nil
}

func (c *CertificateGenerator) PemEncodedCa() []byte {
	return c.CaPem
}

func (c *CertificateGenerator) GenerateKeyAndCertSet(hosts []string, validityInDays int) ([]byte, []byte, error) {
	return c.createKeyPairFromCA(hosts, validityInDays)
}

func (c *CertificateGenerator) generateCa() (*x509.Certificate, *rsa.PrivateKey, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: c.Organizations,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, c.CaValidity),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	c.CaPem = caPEM.Bytes()
	c.CaPrivKeyPem = caPrivKeyPEM.Bytes()

	return ca, caPrivKey, nil
}

func (c *CertificateGenerator) createKeyPairFromCA(hosts []string, validityDays int) ([]byte, []byte, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: c.Organizations,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, validityDays),
		SubjectKeyId:          []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	for _, host := range hosts {
		cert.DNSNames = append(cert.DNSNames, host)
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, c.Ca, &certPrivKey.PublicKey, c.CaPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return certPEM.Bytes(), certPrivKeyPEM.Bytes(), nil
}
