package hippokrates

import (
	"net/http"
	"net/url"
	"path"
)

type odysseiaApi interface {
	Sokrates()
	Herodotos()
	Alexandros()
	Dionysos()
}
type BaseApi struct {
	BaseUrl string
	ApiName string
	Version string
}

type Dionysos struct {
	BaseApi
	Endpoints DionysosEndpoints
}
type Alexandros struct {
	BaseApi
	Endpoints AlexandrosEndpoints
}

type Sokrates struct {
	BaseApi
	Endpoints SokratesEndpoints
}

type Herodotos struct {
	BaseApi
	Endpoints HerodotosEndpoints
}

type CommonEndpoints struct {
	Ping string
	Health string
}

type SokratesEndpoints struct {
	CommonEndpoints
	FindHighestChapter string
	CreateQuestion string
	CheckAnswer string
}

type HerodotosEndpoints struct {
	CommonEndpoints
	CreateSentence string
	CheckAuthor string
}

type AlexandrosEndpoints struct {
	CommonEndpoints
	SearchWord string
}

type DionysosEndpoints struct {
	CommonEndpoints
	CheckGrammar string
}

func (s *Sokrates)GenerateEndpoints() SokratesEndpoints {
	return SokratesEndpoints{
		CommonEndpoints: CommonEndpoints {
			Ping:              "ping",
			Health:             "health",
		},
		FindHighestChapter: "chapters",
		CreateQuestion:     "createQuestion",
		CheckAnswer:        "answer",
	}
}

func (h *Herodotos)GenerateEndpoints() HerodotosEndpoints {
	return HerodotosEndpoints{
		CommonEndpoints: CommonEndpoints {
			Ping:              "ping",
			Health:             "health",
		},
		CreateSentence: "createQuestion",
		CheckAuthor: "authors",
	}
}

func (a *Alexandros)GenerateEndpoints() AlexandrosEndpoints {
	return AlexandrosEndpoints{
		CommonEndpoints: CommonEndpoints {
			Ping:              "ping",
			Health:             "health",
		},
		SearchWord: "search",
	}
}

func (d *Dionysos)GenerateEndpoints() DionysosEndpoints {
	return DionysosEndpoints{
		CommonEndpoints: CommonEndpoints {
			Ping:              "ping",
			Health:             "health",
		},
		CheckGrammar: "checkGrammar",
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

func (h *Herodotos) Health() (*http.Response, error) {
	u, err := url.Parse(h.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, h.ApiName, h.Version, h.Endpoints.Health)

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *Alexandros) Health() (*http.Response, error) {
	u, err := url.Parse(a.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, a.ApiName, a.Version, a.Endpoints.Health)

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}


func (a *Alexandros) QueryWord(word string) (*http.Response, error) {
	u, err := url.Parse(a.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, a.ApiName, a.Version, a.Endpoints.SearchWord)
	q := u.Query()
	q.Set("word", word)
	u.RawQuery = q.Encode()

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

func (h *Herodotos) CreateSentence(author string) (*http.Response, error) {
	u, err := url.Parse(h.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, h.ApiName, h.Version, h.Endpoints.CreateSentence)
	q := u.Query()
	q.Set("author", author)
	u.RawQuery = q.Encode()

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *Dionysos) Health() (*http.Response, error) {
	u, err := url.Parse(d.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, d.ApiName, d.Version, d.Endpoints.Health)

	response, err := GetRequest(*u)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *Dionysos) CheckGrammar(word string) (*http.Response, error) {
	u, err := url.Parse(d.BaseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, d.ApiName, d.Version, d.Endpoints.CheckGrammar)
	q := u.Query()
	q.Set("word", word)
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