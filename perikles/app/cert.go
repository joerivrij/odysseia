package app

import (
	"fmt"
)

func (p *PeriklesHandler) createCert(hosts []string, validityDays int, name, secretName string) error {
	crt, key, err := p.Config.Cert.GenerateKeyAndCertSet(hosts, validityDays)
	if err != nil {
		return err
	}
	certData := make(map[string][]byte)
	certData[fmt.Sprintf("%s.key", name)] = key
	certData[fmt.Sprintf("%s.crt", name)] = crt

	err = p.Config.Kube.Configuration().UpdateSecret(p.Config.Namespace, secretName, certData)
	if err != nil {
		return err
	}
	return nil
}
