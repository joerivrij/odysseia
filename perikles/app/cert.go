package app

import (
	"fmt"
	"github.com/odysseia/plato/generator"
)

func (p *PeriklesHandler) createCert(hosts, organizations []string, name, secretName string) error {
	ca, privateKey, err := generator.GenerateCa(hosts, organizations)
	if err != nil {
		return err
	}

	crt, key, err := generator.CreateKeyPairWithCa(ca, privateKey)
	if err != nil {
		return err
	}
	certData := make(map[string][]byte)
	certData[fmt.Sprintf("%s.key", name)] = key
	certData[fmt.Sprintf("%s.crt", name)] = crt

	err = p.Config.Kube.Configuration().UpdateSecret("odysseia", secretName, certData)
	if err != nil {
		return err
	}
	return nil
}
