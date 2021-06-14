package hippokrates

import (
	"context"
)

const version = "v1"

type LexikoFixture struct {
	ctx		 context.Context
	sokrates Sokrates
	herodotos Herodotos
	alexandros Alexandros
}

func New(alexandrosUrl, herodotosUrl, sokratesUrl, sokratesName, herodotosName, alexandrosName string) (*LexikoFixture, error) {
	sokratesApi := Sokrates{
		BaseApi: BaseApi{
			BaseUrl:   sokratesUrl,
			ApiName:   sokratesName,
			Version: 	version,
		},
		Endpoints: SokratesEndpoints{},
	}
	sokratesApi.Endpoints = sokratesApi.GenerateEndpoints()

	herodotosApi := Herodotos{
		BaseApi: BaseApi{
			BaseUrl:   herodotosUrl,
			ApiName:   herodotosName,
			Version: 	version,
		},
		Endpoints: HerodotosEndpoints{},
	}
	herodotosApi.Endpoints = herodotosApi.GenerateEndpoints()

	alexandrosApi := Alexandros{
		BaseApi: BaseApi{
			BaseUrl:   alexandrosUrl,
			ApiName:   alexandrosName,
			Version: 	version,
		},
		Endpoints: AlexandrosEndpoints{},
	}
	alexandrosApi.Endpoints = alexandrosApi.GenerateEndpoints()

	return &LexikoFixture{
		sokrates:                    sokratesApi,
		herodotos: herodotosApi,
		alexandros: alexandrosApi,
		ctx:                         context.Background(),
	}, nil
}
