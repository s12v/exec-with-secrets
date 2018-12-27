FROM amazonlinux:2

COPY secure-exec-linux-amd64 /usr/local/bin/secure-exec

CMD secure-exec
