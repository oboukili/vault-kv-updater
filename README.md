# Vault KV updater


Simple, no-dependency golang application that flattens and synchronizes a YAML formatted structure from an optionally [SOPS](https://github.com/mozilla/sops) encrypted stream (file or standard input), to a Vault KV secret.

Write operations on secrets are **idempotent**, they will not update secrets if no changes have been detected.

Pretty useful when used in conjunction with [Spring Cloud Vault](https://github.com/spring-cloud/spring-cloud-vault).

Supported Vault authentication methods:
- kubernetes

---

#### Roadmap

* Opt-out flattening flag.
* Opt-out sops features flag.
* Support for simple Vault token authentication mode.
* Support for more flattening modes.
* Providing a useful CLI help.
* Support for multiple files/streams (go routines) + Vault API rate limiting.
* Less hacky yaml marshalling using https://github.com/ghodss/yaml.

#### Usage

##### Simple mode (single secret)
```
vault-kv-updater secret.yml
```

```
cat secret.yml | vault-kv-updater
```

##### Autocomplete mode (multiple secrets)
```
AUTO_COMPLETE=true AUTO_COMPLETE_VAULT_KV_MOUNT=kv vault-kv-updater some/folder
```

```
AUTO_COMPLETE=true AUTO_COMPLETE_VAULT_KV_MOUNT=kv AUTO_COMPLETE_FILE_PREFIX="application-" vault-kv-updater some/folder
```


#### Autocomplete mode

Autocomplete mode will automatically determine the secrets KV destination path based on the filenames.

A few environment variables are provided to give further control, such as adding a KV path destination prefix,
or escaping a common secret name prefix from the filenames.   

Supported file extensions are `.yml` and `.yaml` 

#### Environment variables

|Variable|Optional|Description|defaults|
|---|---|---|---|
|VAULT_ADDR|yes|Vault endpoint address, including scheme and port|"http://127.0.0.1:8200"|
|VAULT_KV_PATH|**no**|Secret path, including kv mount||
|VAULT_ROLE|**no**|Vault role to authenticate against||
|VAULT_CAPEM|yes|Vault CA certificate in PEM format||
|VAULT_CACERT|yes|Path to the vault CA file||
|VAULT_NAMESPACE|yes|Vault namespace (enterprise feature)||
|VAULT_TLS_SERVER_NAME|yes|Vault server hostname to verify against||
|VAULT_SKIP_VERIFY|yes|Whether to skip TLS verification|false|
|VAULT_K8S_MOUNT_PATH|yes|Authentication backend mount path|"kubernetes"|
|SERVICE_ACCOUNT_PATH|yes|Path to the Kubernetes serviceaccount token file|"/var/run/secrets/kubernetes.io/serviceaccount/token"|
|AUTO_COMPLETE|yes|Activates autocomplete mode|false|
|AUTO_COMPLETE_FILE_PREFIX|yes|Removes the prefix from the filename before determining the associated Vault secret's KV path||
|AUTO_COMPLETE_VAULT_KV_MOUNT|yes (autocomplete:**no**)|Vault KV mount to synchronize secrets to||
|AUTO_COMPLETE_VAULT_KV_PATH_PREFIX|yes|Appends a base KV path, i.e. kv/mybasekvpath/secretname||

---

#### Credits

Uses the following OSS libs, thanks guys ;)
* https://github.com/jeremywohl/flatten

Inspired from chunks of the following OSS projects, thanks people :D
* https://github.com/sethvargo/vault-kubernetes-authenticator