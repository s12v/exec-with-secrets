# secure-exec

Wraps an application and opulates secrets from AWS KMS and SSM

## Example:

### Wrap an executable

```
PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." secure-exec app
```

`PARAM` will be decrypted (using default AWS credentials chain) and passed to `app` via environment.

### Docker example

```
FROM amazonlinux:2

ADD https://github.com/secure-exec /secure-exec

CMD /secure-exec java -jar /myapp.jar
```

```
docker run \
    -e PLAINTEXT_PARAM="text" \
    -e KMS_PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." \
    myapp 
```

`KMS_PARAM` will be decrypted and passed to `myapp`. `docker inspect` will still see encrypted value, only `myapp` receives plaintext.
