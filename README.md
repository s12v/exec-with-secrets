[![Build Status](https://travis-ci.com/s12v/secure-exec.svg?branch=master)](https://travis-ci.com/s12v/secure-exec)
[![codecov](https://codecov.io/gh/s12v/secure-exec/branch/master/graph/badge.svg)](https://codecov.io/gh/s12v/secure-exec)

# secure-exec

`secure-exec` populates secrets using AWS KMS or SSM into your app.

It looks for prefixed variables in environment and replaces them:
 - `{aws-kms}encrypted-text` - decrypts the value using AWS KMS
 - `{aws-ssm}parameter-name` - loads parameters from AWS Systems Manager Parameter Store
 
Then it runs `exec` system call and replaces itself with your app.
 
For AWS access the default credentials chain is used. 

## Examples

### Wrap an executable

```
PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." secure-exec app
```

`PARAM` will be decrypted and passed to `app` via environment.

### Docker example

Build an image:

```
FROM amazonlinux:2

ADD https://github.com/secure-exec/ /secure-exec

COPY app.jar /app.jar

CMD secure-exec java -jar /app.jar
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
`docker inspect` will still see the old values
