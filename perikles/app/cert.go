package app

import (
	"fmt"
	"github.com/kpango/glg"
	"time"
)

func (p *PeriklesHandler) createCert(hosts []string, validityDays int, secretName string) error {
	tlsName := "tls"
	crt, key, err := p.Config.Cert.GenerateKeyAndCertSet(hosts, validityDays)
	if err != nil {
		return err
	}

	certData := make(map[string][]byte)
	certData[fmt.Sprintf("%s.key", tlsName)] = key
	certData[fmt.Sprintf("%s.crt", tlsName)] = crt
	certData[fmt.Sprintf("%s.pem", tlsName)] = p.Config.Cert.PemEncodedCa()

	secret, _ := p.Config.Kube.Configuration().GetSecret(p.Config.Namespace, secretName)

	if secret == nil {
		glg.Infof("secret %s does not exist", secretName)
		err = p.Config.Kube.Configuration().CreateTlSSecret(p.Config.Namespace, secretName, certData, false)
		if err != nil {
			return err
		}
	} else {
		glg.Infof("secret %s already exists", secret.Name)

		newAnnotation := make(map[string]string)
		newAnnotation[AnnotationUpdate] = time.Now().UTC().Format(timeFormat)

		err = p.Config.Kube.Configuration().UpdateTLSSecret(p.Config.Namespace, secretName, certData, newAnnotation)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
