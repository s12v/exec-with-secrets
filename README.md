# secure-exec

Populates secrets inside a Docker container from AWS KMS and SSM

Example:

```
FROM amazonlinux:2

ADD https://github.com/secure-exec /secure-exec

CMD /secure-exec java -jar /myapp.jar
```

```
docker run \
    -e PLAINTEXT_PARAM="text" \
    -e KMS_PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." \
    -e SSM_PARAM="{aws-ssm}/myapp/password" \
    myapp 
```

Decrypt secrets and populate application environment
Populate environment with decrypted secrets


