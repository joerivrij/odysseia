package aristoteles

import (
	"github.com/odysseia-greek/plato/certificates"
	"strconv"
)

func (c *Config) getCertClient() (certificates.CertClient, error) {
	envCaValidity := c.getStringFromEnv(EnvCAValidity, defaultCaValidity)
	caValidity, err := strconv.Atoi(envCaValidity)
	if err != nil {
		return nil, err
	}

	org := []string{
		"odysseia",
	}

	certClient, err := certificates.NewCertGeneratorClient(org, caValidity)
	if err != nil {
		return nil, err
	}

	return certClient, nil
}
