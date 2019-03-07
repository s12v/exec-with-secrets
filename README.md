[![Build Status](https://travis-ci.com/s12v/secure-exec.svg?branch=master)](https://travis-ci.com/s12v/secure-exec)
[![codecov](https://codecov.io/gh/s12v/secure-exec/branch/master/graph/badge.svg)](https://codecov.io/gh/s12v/secure-exec)

# secure-exec

`secure-exec` passes secrets from AWS KMS, SSM, or Secrets Manager into your app environment in a secure way.

It supports the following services as secrets providers:
 - [AWS Key Management (KMS)](https://aws.amazon.com/kms/)
 - [AWS Systems Manager Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
 - [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)

This small utility looks for prefixed variables in environment and replaces them with the secret value:
 - `{aws-kms}AQICAHjA3mwbmf...` - decrypts the value using AWS KMS
 - `{aws-ssm}/app/staging/param` - loads parameter `/app/staging/param` from AWS Systems Manager Parameter Store
 - `{aws-sm}/app/staging/param` - loads secret `/app/staging/param` from AWS Secrets Manager
 - `{aws-sm}/app/staging/param{prop1}` - loads secret `/app/staging/param` from AWS Secrets Manager and takes `prop1` property
 
Then it runs `exec` system call and replaces itself with your app.
The secrets are only available to your application and not accessible with `docker inspect`.

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

`KMS_PARAM` and `SSM_PARAM` will be decrypted and passed to `myapp` environment.
`docker inspect` will still see the encrypted values

## Build

`make` builds Linux and Mac binaries with all providers.

To chose providers (for example only AWS SSM), run:
```
make TAGS=awsssm
```

## Adding a new provider

See example PR: https://github.com/s12v/secure-exec/pull/1
