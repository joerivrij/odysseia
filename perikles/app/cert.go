package app

import (
	"fmt"
	"github.com/kpango/glg"
)

func (p *PeriklesHandler) createCert(hosts []string, validityDays int, secretName string) error {
	tlsName := "tls"
	crt, key, err := p.Config.Cert.GenerateKeyAndCertSet(hosts, validityDays)
	if err != nil {
		glg.Error(err)
		return err
	}

	certData := make(map[string][]byte)
	certData[fmt.Sprintf("%s.key", tlsName)] = key
	certData[fmt.Sprintf("%s.crt", tlsName)] = crt

	secret, _ := p.Config.Kube.Configuration().GetSecret(p.Config.Namespace, secretName)

	if secret == nil {
		glg.Info("secret %s does not exist", secretName)
		err = p.Config.Kube.Configuration().CreateTlSSecret(p.Config.Namespace, secretName, certData)
		if err != nil {
			return err
		}
	} else {
		glg.Infof("secret %s already exists", secret.Name)
		err = p.Config.Kube.Configuration().UpdateTLSSecret(p.Config.Namespace, secretName, certData)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
