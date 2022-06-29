module github.com/odysseia

go 1.16

require (
	github.com/elastic/go-elasticsearch/v7 v7.12.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/vault/api/auth/kubernetes v0.1.0
	github.com/kpango/glg v1.5.8
	github.com/kubemq-io/kubemq-go v1.7.2
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/text v0.3.7
	gopkg.in/oauth2.v3 v3.12.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.7.2
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/cli-runtime v0.22.4
	k8s.io/client-go v0.22.4
)

require (
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/hashicorp/vault v1.10.3
	github.com/hashicorp/vault/api v1.5.0
	github.com/michaelklishin/rabbit-hole/v2 v2.12.0 // indirect
	k8s.io/apiextensions-apiserver v0.22.4
)
