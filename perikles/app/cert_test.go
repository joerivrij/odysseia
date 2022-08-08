package app

import (
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/certificates"
	"github.com/odysseia/plato/kubernetes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCertCreation(t *testing.T) {
	ns := "test"
	secretName := "testsecret"
	hosts := []string{
		"perikles",
		"perikles.odysseia",
		"perikles.odysseia.svc",
		"perikles.odysseia.svc.cluster.local",
	}
	organizations := []string{"test"}
	validityCa := 3650

	cert, err := certificates.NewCertGeneratorClient(organizations, validityCa)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	err = cert.InitCa()
	assert.Nil(t, err)

	t.Run("SecretNewlyCreated", func(t *testing.T) {
		fakeKube, err := kubernetes.FakeKubeClient(ns)
		assert.Nil(t, err)
		testConfig := configs.PeriklesConfig{
			Kube:      fakeKube,
			Cert:      cert,
			Namespace: ns,
		}

		handler := PeriklesHandler{Config: &testConfig}
		err = handler.createCert(hosts, 1, secretName)
		assert.Nil(t, err)
	})

	t.Run("SecretAlreadyExists", func(t *testing.T) {
		fakeKube, err := kubernetes.FakeKubeClient(ns)
		assert.Nil(t, err)
		testConfig := configs.PeriklesConfig{
			Kube:      fakeKube,
			Cert:      cert,
			Namespace: ns,
		}

		data := map[string][]byte{
			"somesecret": []byte("verysecret"),
		}

		fakeKube.Configuration().CreateSecret(ns, secretName, data)

		handler := PeriklesHandler{Config: &testConfig}
		err = handler.createCert(hosts, 1, secretName)
		assert.Nil(t, err)
	})
}
