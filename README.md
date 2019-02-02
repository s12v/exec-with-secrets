[![Build Status](https://travis-ci.com/s12v/secure-exec.svg?branch=master)](https://travis-ci.com/s12v/secure-exec)
[![codecov](https://codecov.io/gh/s12v/secure-exec/branch/master/graph/badge.svg)](https://codecov.io/gh/s12v/secure-exec)

# secure-exec

`secure-exec` populates secrets using AWS KMS, SSM, or Secrets Manager into your app environment.

It looks for prefixed variables in environment and replaces them:
 - `{aws-kms}AQICAHjA3mwbmf...` - decrypts the value using AWS KMS
 - `{aws-ssm}/app/staging/param` - loads parameter `/app/staging/param` from AWS Systems Manager Parameter Store
 - `{aws-sm}/app/staging/param` - loads secret `/app/staging/param` from AWS Secrets Manager
 - `{aws-sm}/app/staging/param{prop1}` - loads secret `/app/staging/param` from AWS Secrets Manager and takes `prop1` property
 
Then it runs `exec` system call and replaces itself with your app.
The default credentials chain is used for AWS access.

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

ADD https://github.com/s12v/secure-exec/releases/download/v0.2.1/secure-exec-linux-amd64 /secure-exec

COPY app.jar /app.jar

CMD secure-exec java -jar /app.jar
```

Run:
```
docker run \
    -e PLAINTEXT_PARAM="text" \
    -e KMS_PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." \
    -e SSM_PARAM="{aws-ssm}/myapp/param" \
    myapp 
```

`KMS_PARAM` and `SSM_PARAM` will be decrypted/populated and passed to `myapp` environment.
`docker inspect` will still see the old values

## Build

`make` will build Linux and Mac binaries with all providers.

To chose providers (for example only AWS SSM), run:
```
make TAGS=awsssm
```

## Adding a new provider

See example PR: https://github.com/s12v/secure-exec/pull/1
