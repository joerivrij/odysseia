package aristoteles

import (
	"crypto/tls"
	"errors"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/service"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

func (c *Config) getOdysseiaClient() (service.OdysseiaClient, error) {
	solonUrl := c.getStringFromEnv(EnvSolonService, defaultSolonService)
	tlsEnabled := c.getBoolFromEnv(EnvTlSKey)

	var certBundle service.CertBundle
	var ca []byte

	glg.Debug("getting odysseia client")

	if tlsEnabled {
		glg.Debug("setting up certs because TLS is enabled")
		rootPath := os.Getenv("CERT_ROOT")
		glg.Debugf("rootPath: %s", rootPath)
		dirs, err := ioutil.ReadDir(rootPath)
		if err != nil {
			glg.Error(err)
			return nil, err
		}

		for _, dir := range dirs {
			if dir.IsDir() {
				dirPath := filepath.Join(rootPath, dir.Name())
				glg.Debugf("found directory: %s", dirPath)

				certPath := filepath.Join(dirPath, "tls.crt")
				keyPath := filepath.Join(dirPath, "tls.key")

				if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
					glg.Debugf("cannot get file because it does not exist: %s", certPath)
					continue
				}

				if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
					glg.Debugf("cannot get file because it does not exist: %s", keyPath)
					continue
				}

				loadedCerts, err := tls.LoadX509KeyPair(certPath, keyPath)
				if err != nil {
					glg.Error(err)
					return nil, err
				}

				if ca == nil {
					caPath := filepath.Join(rootPath, dir.Name(), "tls.pem")
					if _, err := os.Stat(caPath); errors.Is(err, os.ErrNotExist) {
						glg.Debugf("cannot get file because it does not exist: %s", caPath)
						continue
					}
					ca, _ = ioutil.ReadFile(caPath)
					glg.Debugf("writing CA for path %s", caPath)
				}

				switch dir.Name() {
				case "solon":
					certBundle.SolonCert = []tls.Certificate{loadedCerts}
				case "ptolemaios":
					certBundle.PtolemaiosCert = []tls.Certificate{loadedCerts}
				case "dionysios":
					certBundle.DionysiosCert = []tls.Certificate{loadedCerts}
				case "herodotos":
					certBundle.HerodotosCert = []tls.Certificate{loadedCerts}
				case "alexandros":
					certBundle.AlexandrosCert = []tls.Certificate{loadedCerts}
				case "sokrates":
					certBundle.SokratesCert = []tls.Certificate{loadedCerts}
				}
			}
		}
	}

	parsedSolon, err := url.Parse(solonUrl)
	if err != nil {
		return nil, err
	}

	config := service.ClientConfig{
		CertBundle:    certBundle,
		Ca:            ca,
		Scheme:        parsedSolon.Scheme,
		SolonUrl:      parsedSolon.Host,
		PtolemaiosUrl: "",
	}

	glg.Debug("creating new client")

	return service.NewClient(config)
}
