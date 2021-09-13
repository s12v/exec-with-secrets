FROM amazonlinux:2

ADD https://github.com/johnrichardrinehart/exec-with-secrets/releases/latest/download/exec-with-secrets-linux-amd64 /exec-with-secrets

RUN chmod +x /exec-with-secrets

ENTRYPOINT ["/exec-with-secrets"]

CMD env
