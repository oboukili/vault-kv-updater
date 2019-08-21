FROM golang:1.12-alpine as builder
ENV GOPATH=/go

RUN apk --no-cache --update add git make
COPY . .
RUN make deps
RUN make build

FROM alpine:3.10.1

COPY --from=builder /go/vault-kv-updater /bin/vault-kv-updater
RUN chmod +x /bin/vault-kv-updater

ENTRYPOINT ["/bin/vault-kv-updater"]