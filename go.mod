module github.com/odysseia

go 1.16

require (
	github.com/elastic/go-elasticsearch/v7 v7.12.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/vault/api/auth/kubernetes v0.1.0
	github.com/kpango/glg v1.5.8
	github.com/kubemq-io/kubemq-go v1.7.2
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.7.2
	golang.org/x/text v0.3.7
	gopkg.in/oauth2.v3 v3.12.0
	gopkg.in/yaml.v3 v3.0.1
	helm.sh/helm/v3 v3.9.4
	k8s.io/api v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/cli-runtime v0.24.2
	k8s.io/client-go v0.24.2
)

require (
	cloud.google.com/go/kms v1.4.0 // indirect
	cloud.google.com/go/monitoring v1.6.0 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/hashicorp/vault v1.10.6
	github.com/hashicorp/vault/api v1.5.0
	k8s.io/apiextensions-apiserver v0.24.2
)
