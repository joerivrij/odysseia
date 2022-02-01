package harbor

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Client interface {
	CreateProject(projectName string, public bool) error
	DeleteProject(projectId string) error
}

type Harbor struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

type ProjectCreationRequest struct {
	ProjectName string   `json:"project_name"`
	Metadata    Metadata `json:"metadata"`
}

type Metadata struct {
	Public string `json:"public"`
}

func NewHarborClient(baseURL, username, password string, cert []byte) (Client, error) {
	var client *http.Client
	if cert == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client = &http.Client{Transport: tr}
	} else {
		caCertPool, _ := x509.SystemCertPool()
		caCertPool.AppendCertsFromPEM(cert)

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}
	}

	harbor := Harbor{
		baseURL:  baseURL,
		username: username,
		password: password,
		client:   client,
	}

	return &harbor, nil
}

func (h *Harbor) CreateProject(projectName string, public bool) error {
	u, err := url.Parse(h.baseURL)
	if err != nil {
		return err
	}

	u.Path = path.Join(u.Path, "/api/v2.0/projects")

	var projectCreationRequest ProjectCreationRequest

	projectCreationRequest.ProjectName = projectName
	if public {
		projectCreationRequest.Metadata.Public = "true"
	} else {
		projectCreationRequest.Metadata.Public = "false"
	}

	d, _ := json.Marshal(&projectCreationRequest)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(d))
	if err != nil {
		return err
	}

	req.SetBasicAuth(h.username, h.password)

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 409 {
		return errors.New("already exists")
	} else if resp.StatusCode != 201 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Could not create project " + projectName + " Body: " + string(respBody) + " status - " + fmt.Sprintf("%v", resp.StatusCode))
	}

	return nil
}

func (h *Harbor) DeleteProject(projectId string) error {
	u, err := url.Parse(h.baseURL)
	if err != nil {
		return err
	}

	urlPath := fmt.Sprintf("api/projects/%s", projectId)
	u.Path = path.Join(u.Path, urlPath)

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(h.username, h.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("got %v were 200 was expected", resp.StatusCode)
	}

	return nil
}
