package certificates

import (
	"github.com/kpango/glg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneration(t *testing.T) {
	hosts := []string{
		"perikles",
		"perikles.odysseia",
		"perikles.odysseia.svc",
		"perikles.odysseia.svc.cluster.local",
	}
	organizations := []string{"test"}
	validityCa := 3650
	validityCert := 3650

	t.Run("Authority", func(t *testing.T) {
		impl, err := NewCertGeneratorClient(organizations, validityCa)
		assert.Nil(t, err)
		assert.NotNil(t, impl)
		err = impl.InitCa()
		assert.Nil(t, err)

		crt, key, err := impl.GenerateKeyAndCertSet(hosts, validityCert)
		glg.Info(string(crt))
		glg.Info(string(key))
		assert.Nil(t, err)
		assert.NotNil(t, key)
		assert.NotNil(t, crt)
	})

}
