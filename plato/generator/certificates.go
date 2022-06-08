package generator

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"time"
)

func GenerateKeyAndCertSet(hosts, organization []string) ([]byte, []byte, error) {
	ca, privateKey, err := generateCa(hosts, organization)
	if err != nil {
		return nil, nil, err
	}

	return generateKeyAndCert(ca, privateKey)
}

func GenerateCa(hosts, organization []string) ([]byte, []byte, error) {
	return generateCa(hosts, organization)
}

func CreateKeyPairWithCa(ca, key []byte) ([]byte, []byte, error) {
	return generateKeyAndCert(ca, key)
}

func generateCa(hosts, organization []string) ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: organization,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, 3650),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	for _, host := range hosts {
		ca.DNSNames = append(ca.DNSNames, host)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(key)

	return caBytes, keyBytes, nil
}

func generateKeyAndCert(caBytes, key []byte) ([]byte, []byte, error) {
	caPEM := new(bytes.Buffer)
	err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	if err != nil {
		return nil, nil, err
	}

	keyBytes, _ := x509.ParsePKCS1PrivateKey(key)
	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyBytes),
	})

	if err != nil {
		return nil, nil, err
	}

	return caPEM.Bytes(), caPrivKeyPEM.Bytes(), nil
}

func whatever(key []byte) ([]byte, error) {
	subj := pkix.Name{
		CommonName:         "example.com",
		Country:            []string{"AU"},
		Province:           []string{"Some-State"},
		Locality:           []string{"MyCity"},
		Organization:       []string{"Company Ltd"},
		OrganizationalUnit: []string{"IT"},
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type: asn1.ObjectIdentifier{},
				Value: asn1.RawValue{
					Tag:   asn1.TagIA5String,
					Bytes: []byte("sdf"),
				},
			},
		},
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	keyBytes, _ := x509.ParsePKCS1PrivateKey(key)
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)

	return csrBytes, err
}
