package hippokrates

import (
	"context"
)

const version = "v1"

type SokratesFixture struct {
	ctx		 context.Context
	sokrates Sokrates
}

func New(baseUrl, apiName string) (*SokratesFixture, error) {

	sokratesApi := Sokrates{
		BaseUrl:   baseUrl,
		ApiName:   apiName,
		Version: 	version,
		Endpoints: GenerateEndpoints(),
	}

	return &SokratesFixture{
		sokrates:                    sokratesApi,
		ctx:                         context.Background(),
	}, nil
}
