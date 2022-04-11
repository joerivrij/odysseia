package hippokrates

import (
	"context"
	"github.com/kpango/glg"
	"github.com/odysseia/hippokrates/client"
	"net/url"
	"os"
)

const (
	defaultAlexandrosService = "http://odysseia-greek.internal"
	defaultDionysosService   = "http://odysseia-greek.internal"
	defaultHerodotosService  = "http://odysseia-greek.internal"
	defaultSokratesService   = "http://odysseia-greek.internal"
	defaultSolonService      = "http://odysseia-greek.internal"
	EnvDionysosService       = "DIONYSOS_SERVICE"
	EnvAlexandrosService     = "ALEXANDROS_SERVICE"
	EnvHerodotosService      = "HERODOTOS_SERVICE"
	EnvSokratesService       = "SOKRATES_SERVICE"
	EnvSolonService          = "SOLON_SERVICE"
)

type odysseiaFixture struct {
	ctx     context.Context
	clients client.OdysseiaClient
}

func New(config *client.ClientConfig) (*odysseiaFixture, error) {
	cfg, err := client.NewClient(*config)

	if err != nil {
		return nil, err
	}

	return &odysseiaFixture{
		clients: cfg,
		ctx:     context.Background(),
	}, nil
}

func GetEnv() (*client.ClientConfig, error) {
	alexandrosUrl := getStringFromEnv(EnvAlexandrosService, defaultAlexandrosService)

	parsedAlexandros, err := url.Parse(alexandrosUrl)
	if err != nil {
		return nil, err
	}

	dionysosUrl := getStringFromEnv(EnvDionysosService, defaultDionysosService)

	parsedDionysos, err := url.Parse(dionysosUrl)
	if err != nil {
		return nil, err
	}

	herodotosUrl := getStringFromEnv(EnvHerodotosService, defaultHerodotosService)

	parsedHerodotos, err := url.Parse(herodotosUrl)
	if err != nil {
		return nil, err
	}

	sokratesUrl := getStringFromEnv(EnvSokratesService, defaultSokratesService)

	parsedSokrates, err := url.Parse(sokratesUrl)
	if err != nil {
		return nil, err
	}

	solonUrl := getStringFromEnv(EnvSolonService, defaultSolonService)

	parsedSolon, err := url.Parse(solonUrl)
	if err != nil {
		return nil, err
	}

	config := client.ClientConfig{
		Scheme:        parsedAlexandros.Scheme,
		AlexandrosUrl: parsedAlexandros.Host,
		DionysosUrl:   parsedDionysos.Host,
		HerodotosUrl:  parsedHerodotos.Host,
		SokratesUrl:   parsedSokrates.Host,
		SolonUrl:      parsedSolon.Host,
	}

	return &config, nil

}

func getStringFromEnv(env, defaultValue string) string {
	var value string
	value = os.Getenv(env)
	if value == "" {
		glg.Debugf("%s empty set as env variable - defaulting to %s", env, defaultValue)
		value = defaultValue
	}

	return value
}
