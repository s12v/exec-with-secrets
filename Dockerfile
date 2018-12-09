FROM amazonlinux:2

COPY secure-exec /usr/local/bin/secure-exec

CMD secure-exec
