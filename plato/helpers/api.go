package helpers

import (
	"bytes"
	"net/http"
	"net/url"
)

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

func PostRequest(u url.URL, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
