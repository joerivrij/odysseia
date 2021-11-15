package vault

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetClusterKeys(namespace string) models.ClusterKeys {
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == "odysseia" {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}
	clusterKeyName := fmt.Sprintf("cluster-keys-%s.json", namespace)
	clusterKeyFilePath := filepath.Join(l, "solon", "vault_config", clusterKeyName)
	f, err := ioutil.ReadFile(clusterKeyFilePath)
	if err != nil {
		glg.Fatal(err)
	}

	clusterKeys, err := models.UnmarshalClusterKeys(f)
	if err != nil {
		glg.Fatal(err)
	}

	return clusterKeys
}

