package main

import (
	"crypto/tls"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/perikles/app"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	standardPort = "4443"
	crtFileName  = "tls.crt"
	keyFileName  = "tls.key"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PERIKLES
	glg.Info("\n ____   ___  ____   ____  __  _  _        ___  _____\n|    \\ /  _]|    \\ |    ||  |/ ]| |      /  _]/ ___/\n|  o  )  [_ |  D  ) |  | |  ' / | |     /  [_(   \\_ \n|   _/    _]|    /  |  | |    \\ | |___ |    _]\\__  |\n|  | |   [_ |    \\  |  | |     ||     ||   [_ /  \\ |\n|  | |     ||  .  \\ |  | |  .  ||     ||     |\\    |\n|__| |_____||__|\\_||____||__|\\_||_____||_____| \\___|\n                                                    \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"τόν γε σοφώτατον οὐχ ἁμαρτήσεται σύμβουλον ἀναμείνας χρόνον.\"")
	glg.Info("\"he would yet do full well to wait for that wisest of all counsellors, Time.\"")
	glg.Info(strings.Repeat("~", 37))

	glg.Debug("creating config")

	baseConfig := configs.PeriklesConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	periklesConfig, ok := unparsedConfig.(*configs.PeriklesConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	handler := app.PeriklesHandler{Config: periklesConfig}

	glg.Info("init for CA started...")
	err = handler.Config.Cert.InitCa()
	if err != nil {
		glg.Fatal(err)
	}

	glg.Info("CA created")

	glg.Info("creating CRD...")
	created, err := handler.Config.Kube.V1Alpha1().ServiceMapping().CreateInCluster()
	if err != nil {
		glg.Error(err)
	}

	if created {
		glg.Info("CRD created")
	} else {
		glg.Info("CRD not created, it might already exist")
	}

	_, err = handler.Config.Kube.V1Alpha1().ServiceMapping().Get(periklesConfig.CrdName)
	if err != nil {
		glg.Error(err)
		mapping, err := handler.Config.Kube.V1Alpha1().ServiceMapping().Parse(nil, periklesConfig.CrdName, periklesConfig.Namespace)
		if err != nil {
			glg.Error(err)
		}

		createdCrd, err := handler.Config.Kube.V1Alpha1().ServiceMapping().Create(mapping)
		if err != nil {
			glg.Error(err)
		}

		glg.Debugf("created mapping %s", createdCrd.Name)

	}

	glg.Debug("init routes")
	srv := app.InitRoutes(*periklesConfig)

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
	}

	glg.Debug("setting up server with https")

	httpsServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      srv,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	glg.Debug("loading cert files from mount")
	certFile := filepath.Join(periklesConfig.TLSFiles, crtFileName)
	keyFile := filepath.Join(periklesConfig.TLSFiles, keyFileName)

	err = httpsServer.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		glg.Fatal(err)
	}
}
