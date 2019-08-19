FROM alpine:3.10.1

COPY /go/vault-kv-updater /bin/vault-kv-updater
RUN chmod +x /bin/vault-kv-updater

ENTRYPOINT ["/bin/vault-kv-updater"]