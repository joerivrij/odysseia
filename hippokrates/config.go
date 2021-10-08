package hippokrates

import (
	"context"
)

const version = "v1"

type odysseiaFixture struct {
	ctx		 context.Context
	sokrates Sokrates
	herodotos Herodotos
	alexandros Alexandros
	dionysos Dionysos
}

func New(alexandrosUrl, herodotosUrl, sokratesUrl, dionysosUrl, sokratesName, herodotosName, alexandrosName, dionysosName string) (*odysseiaFixture, error) {
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

	dionysosApi := Dionysos {
		BaseApi: BaseApi{
			BaseUrl:   dionysosUrl,
			ApiName:   dionysosName,
			Version: 	version,
		},
		Endpoints: DionysosEndpoints{},
	}
	dionysosApi.Endpoints = dionysosApi.GenerateEndpoints()

	return &odysseiaFixture{
		sokrates:                    sokratesApi,
		herodotos: herodotosApi,
		alexandros: alexandrosApi,
		dionysos: dionysosApi,
		ctx:                         context.Background(),
	}, nil
}
