package elastic

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	fixtures = make(map[string]io.ReadCloser)
)

func init() {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == "odysseia" {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}
	eratosthenesDir := filepath.Join(l, "eratosthenes", "*.json")
	fixtureFiles, err := filepath.Glob(eratosthenesDir)
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
	case 500:
		mockCode = http.StatusInternalServerError
	case 502:
		mockCode = http.StatusBadGateway
	default:
		mockCode = 200
	}

	body := fixture(fmt.Sprintf("%s.json", fixtureFile))

	mockTrans := MockTransport{
		Response: &http.Response{
			StatusCode: mockCode,
			Body:       body,
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

func CreateEmptyClient() (*elasticsearch.Client, error) {
	mockTrans := MockTransport{
		Response: &http.Response{
			StatusCode: 503,
			Body:       nil,
		},
	}
	mockTrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mockTrans.Response, nil }

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}
