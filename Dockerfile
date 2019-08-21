FROM golang:1.12-alpine as builder
ENV GOPATH=/go

COPY . .
RUN apk --no-cache --update add git make
RUN make deps
RUN make build

FROM alpine:3.10.1

COPY --from=builder /go/vault-kv-updater /bin/vault-kv-updater
RUN chmod +x /bin/vault-kv-updater

ENTRYPOINT ["/bin/vault-kv-updater"]
