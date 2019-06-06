FROM amazonlinux:2

ADD https://github.com/s12v/exec-with-secrets/releases/download/v0.4.0/exec-with-secrets-linux-amd64 /exec-with-secrets

RUN chmod +x /exec-with-secrets

ENTRYPOINT ["/exec-with-secrets"]

CMD env
