# Vault KV updater


Simple, no-dependency golang application that flattens and synchronizes a YAML formatted structure from an optionally [SOPS](https://github.com/mozilla/sops) encrypted stream (file or standard input), to a Vault KV secret.

Pretty useful when used in conjunction with [Spring Cloud Vault](https://github.com/spring-cloud/spring-cloud-vault).

Supported Vault authentication methods:
- kubernetes

---

##### Roadmap

* Opt-out flattening flag.
* Opt-out sops features flag.
* Support for simple Vault token authentication mode.
* Support for more flattening modes.
* Providing a useful CLI help.
* Support for multiple files/streams (go routines) + Vault API rate limiting.
* Less hacky yaml marshalling using https://github.com/ghodss/yaml.

##### Usage

```
vault-kv-updater secret.yml
```

```
cat secret.yml | vault-kv-updater
```

##### Environment variables

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

---

##### Credits

Uses the following OSS libs, thanks guys ;)
* https://github.com/jeremywohl/flatten

Inspired from chunks of the following OSS projects, thanks people :D
* https://github.com/sethvargo/vault-kubernetes-authenticator