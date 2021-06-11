package hippokrates

import (
	"net/http"
	"net/url"
	"path"
)

type BaseUrl interface {
}

type Sokrates struct {
	BaseUrl string
	ApiName string
	Version string
	Endpoints SokratesEndpoints
}

type SokratesEndpoints struct {
	Ping string
	Health string
	FindHighestChapter string
	CreateQuestion string
	CheckAnswer string
}

func GenerateEndpoints() SokratesEndpoints {
	return SokratesEndpoints{
		Ping:               "ping",
		Health:             "health",
		FindHighestChapter: "chapters",
		CreateQuestion:     "createQuestion",
		CheckAnswer:        "answer",
	}
}
func (s *Sokrates) Ping() (*http.Response, error) {
	u, err := url.Parse(s.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, s.ApiName, s.Version, s.Endpoints.Ping)

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Sokrates) Health() (*http.Response, error) {
	u, err := url.Parse(s.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, s.ApiName, s.Version, s.Endpoints.Health)

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Sokrates) CreateQuestion(category, chapter string) (*http.Response, error) {
	u, err := url.Parse(s.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, s.ApiName, s.Version, s.Endpoints.CreateQuestion)
	q := u.Query()
	q.Set("category", category)
	q.Add("chapter", chapter)
	u.RawQuery = q.Encode()

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetRequest(u url.URL) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}