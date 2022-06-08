package generator

import (
	"github.com/kpango/glg"
	"github.com/odysseia/plato/kubernetes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneration(t *testing.T) {
	t.Run("Authority", func(t *testing.T) {
		hosts := []string{
			"vault",
			"vault.odysseia",
			"vault.odysseia.svc",
			"vault.odysseia.svc.cluster.local",
		}
		organizations := []string{"test"}
		crt, key, err := GenerateKeyAndCertSet(hosts, organizations)
		assert.Nil(t, err)

		certData := make(map[string][]byte)
		certData["vault.key"] = key
		certData["vault.crt"] = crt

		secretName := "vault-server-tls"

		homeDir, _ := os.UserHomeDir()
		kubePath := filepath.Join(homeDir, "/.kube/config")
		cfg, _ := ioutil.ReadFile(kubePath)
		kube, _ := kubernetes.NewKubeClient(cfg, "test")
		err = kube.Configuration().CreateSecret("test", secretName, certData)
		assert.Nil(t, err)
		assert.NotNil(t, key)
		assert.NotNil(t, crt)
	})

	t.Run("Authority", func(t *testing.T) {
		hosts := []string{
			"vault",
			"vault.odysseia",
			"vault.odysseia.svc",
			"vault.odysseia.svc.cluster.local",
		}
		organizations := []string{"test"}

		ca, privateKey, err := GenerateCa(hosts, organizations)
		assert.Nil(t, err)

		crt, key, err := CreateKeyPairWithCa(ca, privateKey)
		glg.Info(string(crt))
		glg.Info(string(key))
		assert.Nil(t, err)
		assert.NotNil(t, key)
		assert.NotNil(t, crt)
	})

}
