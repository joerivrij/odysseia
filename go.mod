module github.com/odysseia

go 1.16

require (
	github.com/elastic/go-elasticsearch/v7 v7.12.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/vault/api v1.3.0
	github.com/hashicorp/vault/api/auth/kubernetes v0.1.0
	github.com/kpango/glg v1.5.8
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	golang.org/x/text v0.3.6
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/oauth2.v3 v3.12.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
)
