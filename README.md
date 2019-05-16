[![Build Status](https://travis-ci.com/s12v/exec-with-secrets.svg?branch=master)](https://travis-ci.com/s12v/exec-with-secrets)
[![codecov](https://codecov.io/gh/s12v/exec-with-secrets/branch/master/graph/badge.svg)](https://codecov.io/gh/s12v/exec-with-secrets)

Populate secrets from AWS KMS, SSM or Secrets Manager into your app environment

`exec-with-secrets` passes secrets from AWS KMS, SSM, or Secrets Manager into your app environment in a secure way.

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
PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." exec-with-secrets app
```

`PARAM` will be decrypted and passed to `app` via environment.

### Docker example

Build an image:

```
FROM amazonlinux:2

ADD https://github.com/s12v/exec-with-secrets/releases/download/v0.2.1/exec-with-secrets-linux-amd64 /exec-with-secrets

COPY app.jar /app.jar

CMD exec-with-secrets java -jar /app.jar
```

Run:
```
docker run \
    -e PLAINTEXT_PARAM="text" \
    -e KMS_PARAM="{aws-kms}AQICAHjA3mwvsfng346vnbmf..." \
    -e SSM_PARAM="{aws-ssm}/myapp/param" \
    myappimage
```

`KMS_PARAM` and `SSM_PARAM` will be decrypted and passed to `app.jar` environment.
`docker inspect` will still see the encrypted values

## Build

`make` builds Linux and Mac binaries with all providers.

To chose providers (for example only AWS SSM), run:
```
make TAGS=awsssm
```

## Adding a new provider

See example PR: https://github.com/s12v/exec-with-secrets/pull/1
