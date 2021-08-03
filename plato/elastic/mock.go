package elastic

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var (
	fixtures = make(map[string]io.ReadCloser)
)

func init() {
	fixtureFiles, err := filepath.Glob("fixtures/*.json")
	if err != nil {
		panic(fmt.Sprintf("Cannot glob fixture files: %s", err))
	}

	for _, fpath := range fixtureFiles {
		f, err := ioutil.ReadFile(fpath)
		if err != nil {
			panic(fmt.Sprintf("Cannot read fixture file: %s", err))
		}
		fixtures[filepath.Base(fpath)] = ioutil.NopCloser(bytes.NewReader(f))
	}
}

func fixture(fname string) io.ReadCloser {
	out := new(bytes.Buffer)
	b1 := bytes.NewBuffer([]byte{})
	b2 := bytes.NewBuffer([]byte{})
	tr := io.TeeReader(fixtures[fname], b1)

	defer func() { fixtures[fname] = ioutil.NopCloser(b1) }()
	io.Copy(b2, tr)
	out.ReadFrom(b2)

	return ioutil.NopCloser(out)
}

type MockTransport struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func CreateMockClient(fixtureFile string, statusCode int) (*elasticsearch.Client, error) {
	mockCode := 500
	switch statusCode {
	case 200:
		mockCode = http.StatusOK
	case 404:
		mockCode = http.StatusNotFound
	case 502:
		mockCode = http.StatusBadGateway
	default:
		mockCode = 200
	}
	mockTrans := MockTransport{
		Response: &http.Response{
			StatusCode: mockCode,
			Body:       fixture(fmt.Sprintf("%s.json", fixtureFile)),
		},
	}
	mockTrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mockTrans.Response, nil }

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Transport: &mockTrans,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}