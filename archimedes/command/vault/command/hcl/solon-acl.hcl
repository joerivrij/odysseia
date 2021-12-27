path "sys/health"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage tokens broadly across Vault
path "auth/token/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete key/value secrets for configs
path "secret/configs/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List, create, update, and delete key/value secrets for configs
path "configs/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}