FROM golang:1.12.9-alpine3.10 as builder
ENV GOPATH=/go

COPY . .
RUN apk --no-cache --update add git ca-certificates tzdata
RUN update-ca-certificates
RUN adduser -D -g '' app

# As long as this repo is private, GITLAB_CREDS should be of the USERNAME:ACCESS_TOKEN form
ARG GITLAB_CREDS
RUN git config --global credential.helper store
RUN echo "https://${GITLAB_CREDS}@gitlab.com" >> ~/.git-credentials && \
    go get -v -d && \
    shred -u ~/.git-credentials

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o vault-kv-updater

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/vault-kv-updater /bin/vault-kv-updater

USER app

ENTRYPOINT ["/bin/vault-kv-updater"]
