package aristoteles

const (
	defaultSidecarService    = "http://127.0.0.1:5001"
	defaultKubeConfig        = "/.kube/config"
	defaultNamespace         = "odysseia"
	defaultPodName           = "somepod-08595-383"
	defaultSearchWord        = "greek"
	defaultRoleName          = "solon"
	defaultChannelName       = "dictionary-channel"
	defaultMqAddress         = "localhost"
	defaultMqPort            = "50000"
	defaultJobName           = "demokritos"
	defaultSolonService      = "http://odysseia-greek.internal"
	defaultCaValidity        = "3650"
	EnvHealthCheckOverwrite  = "HEALTH_CHECK_OVERWRITE"
	EnvPodName               = "POD_NAME"
	EnvNamespace             = "NAMESPACE"
	EnvIndex                 = "ELASTIC_ACCESS"
	EnvSecondaryIndex        = "ELASTIC_SECONDARY_ACCESS"
	EnvVaultService          = "VAULT_SERVICE"
	EnvSolonService          = "SOLON_SERVICE"
	EnvPtolemaiosService     = "PTOLEMAIOS_SERVICE"
	EnvRunOnce               = "RUN_ONCE"
	EnvTlSKey                = "TLS_ENABLED"
	EnvKey                   = "ENV"
	EnvSearchWord            = "SEARCH_WORD"
	EnvRole                  = "ELASTIC_ROLE"
	EnvRoles                 = "ELASTIC_ROLES"
	EnvIndexes               = "ELASTIC_INDEXES"
	EnvRootToken             = "VAULT_ROOT_TOKEN"
	EnvAuthMethod            = "AUTH_METHOD"
	EnvVaultRole             = "VAULT_ROLE"
	EnvKubePath              = "KUBE_PATH"
	EnvSidecarOverwrite      = "SIDECAR_OVERWRITE"
	EnvChannel               = "CHANNEL"
	EnvMqAddress             = "MQ_SERVICE"
	EnvMqPort                = "MQ_PORT"
	EnvJobName               = "JOB_NAME"
	EnvCAValidity            = "CA_VALIDITY"
	AuthMethodKube           = "kubernetes"
	AuthMethodToken          = "token"
	baseDir                  = "base"
	configFileName           = "config.yaml"
	defaultRoleAnnotation    = "odysseia-greek/role"
	defaultAccessAnnotation  = "odysseia-greek/access"
	serviceAccountTokenPath  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	certPathInPod            = "/app/config/certs/elastic-certificate.pem"
	elasticServiceDefault    = "http://localhost:9200"
	elasticServiceDefaultTlS = "https://localhost:9200"
	elasticUsernameDefault   = "elastic"
	elasticPasswordDefault   = "odysseia"
	EnvElasticService        = "ELASTIC_SEARCH_SERVICE"
	EnvElasticUser           = "ELASTIC_SEARCH_USER"
	EnvElasticPassword       = "ELASTIC_SEARCH_PASSWORD"
)

var serviceMapping = map[string]string{
	"SolonService": EnvSolonService,
}

var validFields = []string{
	"SolonService",
	"Index",
	"SecondaryIndex",
	"Created",
	"SearchWord",
	"FullPodName",
	"VaultService",
	"Kube",
	"Elastic",
	"PodName",
	"RunOnce",
	"Namespace",
	"DeclensionConfig",
	"Roles",
	"Indexes",
	"SolonCreationRequest",
	"Vault",
	"ElasticCert",
	"AccessAnnotation",
	"RoleAnnotation",
	"Channel",
	"Job",
	"Queue",
	"HttpClients",
	"Cache",
	"Cert",
}
