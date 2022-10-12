package service

import (
	"crypto/tls"
	"github.com/odysseia/plato/models"
)

type OdysseiaClient interface {
	Solon() Solon
	Ptolemaios() Ptolemaios
}

type Odysseia struct {
	solon      *SolonImpl
	ptolemaios *PtolemaiosImpl
}

type Solon interface {
	Health() (*models.Health, error)
	OneTimeToken() (*models.TokenResponse, error)
	Register(requestBody models.SolonCreationRequest) (*models.SolonResponse, error)
}

type Ptolemaios interface {
	GetSecret() (*models.ElasticConfigVault, error)
}

type ClientConfig struct {
	Ca            []byte
	CertBundle    CertBundle
	Scheme        string
	SolonUrl      string
	PtolemaiosUrl string
}

type CertBundle struct {
	SolonCert      []tls.Certificate
	PtolemaiosCert []tls.Certificate
	DionysiosCert  []tls.Certificate
	HerodotosCert  []tls.Certificate
	AlexandrosCert []tls.Certificate
	SokratesCert   []tls.Certificate
}

func NewClient(config ClientConfig) (OdysseiaClient, error) {
	solonImpl, err := NewSolonImpl(config.Scheme, config.SolonUrl, config.Ca, config.CertBundle.SolonCert)
	if err != nil {
		return nil, err
	}

	ptolemaiosImpl, err := NewPtolemaiosConfig(config.Scheme, config.PtolemaiosUrl, config.Ca, config.CertBundle.PtolemaiosCert)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		ptolemaios: ptolemaiosImpl,
	}, nil
}

func NewFakeClient(config ClientConfig, codes []int, responses []string) (OdysseiaClient, error) {
	client := NewFakeHttpClient(responses, codes)

	solonImpl, err := NewFakeSolonImpl(config.Scheme, config.SolonUrl, client)
	if err != nil {
		return nil, err
	}

	ptolemaiosImpl, err := NewFakePtolemaiosConfig(config.Scheme, config.PtolemaiosUrl, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{
		solon:      solonImpl,
		ptolemaios: ptolemaiosImpl,
	}, nil
}

func (o *Odysseia) Solon() Solon {
	if o == nil {
		return nil
	}
	return o.solon
}

func (o *Odysseia) Ptolemaios() Ptolemaios {
	if o == nil {
		return nil
	}
	return o.ptolemaios
}
