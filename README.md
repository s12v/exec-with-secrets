# secure-exec

Populates secrets using AWS KMS or SSM into your app

## Examples

`secure-exec` looks for prefixed variables in environment and replaces them with secret values:
 - `aws-kms` - decrypts the value using default AWS credentials chain
 - `aws-ssm` - loads parameters from AWS Systems Manager Parameter Store

### Wrap an executable

```
PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." secure-exec app
```

`PARAM` will be decrypted and passed to `app` via environment.

### Docker example

Build an image

```
FROM amazonlinux:2

ADD https://github.com/secure-exec /secure-exec

CMD /secure-exec java -jar /myapp.jar
```

Run:
```
docker run \
    -e PLAINTEXT_PARAM="text" \
    -e KMS_PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." \
    -e SSM_PARAM="{aws-kms}/myapp/param" \
    myapp 
```

`KMS_PARAM` and `SSM_PARAM` will be decrypted/populated and passed to `myapp` environment.
`docker inspect` will still see encrypted value, only `myapp` receives plaintext.
