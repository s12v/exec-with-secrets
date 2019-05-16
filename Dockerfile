FROM amazonlinux:2

COPY ./bin/exec-with-secrets-linux-amd64 /usr/local/bin/exec-with-secrets

CMD exec-with-secrets
