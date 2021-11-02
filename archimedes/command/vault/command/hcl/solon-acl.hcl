path "sys/health"
{
  capabilities = ["read", "sudo"]
}

# Manage tokens broadly across Vault
path "auth/token/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete key/value secrets for configs
path "secret/config/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}