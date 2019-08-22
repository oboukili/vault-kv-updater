package api

const (
	EnvAutoComplete                  = "AUTO_COMPLETE"
	EnvAutoCompleteFilePrefix        = "AUTO_COMPLETE_FILE_PREFIX"
	EnvAutoCompleteVaultKVPathPrefix = "AUTO_COMPLETE_VAULT_KV_PATH_PREFIX"
	EnvVaultKvMount                  = "VAULT_KV_MOUNT"
	EnvVaultKvPath                   = "VAULT_KV_PATH"
	EnvVaultKvVersion                = "VAULT_KV_VERSION"
	EnvVaultAddr                     = "VAULT_ADDR"
	EnvVaultToken                    = "VAULT_TOKEN"
	EnvVaultAuthMethod               = "VAULT_AUTH_METHOD"
	EnvVaultAuthKubernetesRole       = "VAULT_AUTH_K8S_ROLE"
	vaultDefaultAddr                 = "http://127.0.0.1:8200"
	vaultDefaultAuthenticationMethod = "kubernetes"
)
