package service

import "github.com/odysseia/plato/models"

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
	Scheme        string
	SolonUrl      string
	PtolemaiosUrl string
}

func NewClient(config ClientConfig) (OdysseiaClient, error) {
	client := NewHttpClient()

	solonImpl, err := NewSolonImpl(config.Scheme, config.SolonUrl, client)
	if err != nil {
		return nil, err
	}

	ptolemaiosImpl, err := NewPtolemaiosConfig(config.Scheme, config.PtolemaiosUrl, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{solon: solonImpl, ptolemaios: ptolemaiosImpl}, nil
}

func NewFakeClient(config ClientConfig, codes []int, responses []string) (OdysseiaClient, error) {
	client := NewFakeHttpClient(responses, codes)

	solonImpl, err := NewSolonImpl(config.Scheme, config.SolonUrl, client)
	if err != nil {
		return nil, err
	}

	ptolemaiosImpl, err := NewPtolemaiosConfig(config.Scheme, config.PtolemaiosUrl, client)
	if err != nil {
		return nil, err
	}

	return &Odysseia{solon: solonImpl, ptolemaios: ptolemaiosImpl}, nil
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
