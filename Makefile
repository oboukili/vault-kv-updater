deps:
	go get -v -d
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-kv-updater
