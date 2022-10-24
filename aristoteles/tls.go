package aristoteles

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/odysseia-greek/plato/helpers"
	"net/http"
	"os"
	"path/filepath"
)

func CreateTlSConfig(port string, ca *x509.CertPool, server *mux.Router) *http.Server {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		//ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: ca,
	}

	return &http.Server{
		Addr:         port,
		Handler:      server,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}

func RetrieveCertPathLocally(testOverwrite bool, service string) (cert string, key string) {
	keyName := "tls.key"
	certName := "tls.crt"

	if testOverwrite {
		glg.Info("trying to read cert file from file")
		rootPath := helpers.OdysseiaRootPath()
		if service == "" {
			service = "solon"
		}
		cert = filepath.Join(rootPath, "eratosthenes", "fixture", service, certName)
		key = filepath.Join(rootPath, "eratosthenes", "fixture", service, keyName)

		return
	} else {
		rootPath := os.Getenv("CERT_ROOT")
		cert = filepath.Join(rootPath, service, certName)
		key = filepath.Join(rootPath, service, keyName)

		glg.Debugf("found certpath: %s - found keypath: %s", cert, key)
	}

	return
}
