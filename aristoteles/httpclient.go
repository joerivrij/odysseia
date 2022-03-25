package aristoteles

import (
	"github.com/odysseia/plato/service"
	"net/url"
)

func (c *Config) getOdysseiaClient() (service.OdysseiaClient, error) {
	solonUrl := c.getStringFromEnv(EnvSolonService, defaultSolonService)

	parsedSolon, err := url.Parse(solonUrl)
	if err != nil {
		return nil, err
	}

	config := service.ClientConfig{
		Scheme:        parsedSolon.Scheme,
		SolonUrl:      parsedSolon.Host,
		PtolemaiosUrl: "",
	}

	return service.NewClient(config)
}
