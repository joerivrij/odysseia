package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient interface {
	Get(u *url.URL) (*http.Response, error)
	Post(u *url.URL, body []byte) (*http.Response, error)
}

type ClientImpl struct {
	client *http.Client
}

type FakeClientImpl struct {
	responseBodies []string
	codes          []int
	index          int
}

func NewHttpClient(caCert []byte, certs []tls.Certificate) HttpClient {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	if caCert != nil {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: certs,
				},
			},
			Timeout: 10 * time.Second,
		}
	}

	return &ClientImpl{
		client: &client,
	}
}

func NewFakeHttpClient(responseBodies []string, codes []int) HttpClient {
	return &FakeClientImpl{
		responseBodies: responseBodies,
		codes:          codes,
		index:          0,
	}
}

func (c *ClientImpl) Get(u *url.URL) (*http.Response, error) {
	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ClientImpl) Post(u *url.URL, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *FakeClientImpl) Get(u *url.URL) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	responseBody := ioutil.NopCloser(strings.NewReader(f.responseBodies[f.index]))

	response := http.Response{
		StatusCode: f.codes[f.index],
		Body:       responseBody,
	}

	if f.index != len(f.codes)-1 {
		f.index++
	}

	return &response, nil
}

func (f *FakeClientImpl) Post(u *url.URL, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	responseBody := ioutil.NopCloser(strings.NewReader(f.responseBodies[f.index]))

	response := http.Response{
		StatusCode: f.codes[f.index],
		Body:       responseBody,
	}

	if f.index != len(f.codes)-1 {
		f.index++
	}

	return &response, nil
}
