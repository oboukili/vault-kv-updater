deps:
	go get \
	  github.com/go-yaml/yaml \
	  github.com/jeremywohl/flatten \
	  go.mozilla.org/sops/decrypt \
	  github.com/hashicorp/vault/api \
	  golang.org/x/net/http2

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-kv-updater
