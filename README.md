[![Build Status](https://travis-ci.com/s12v/exec-with-secrets.svg?branch=master)](https://travis-ci.com/s12v/exec-with-secrets)
[![codecov](https://codecov.io/gh/s12v/exec-with-secrets/branch/master/graph/badge.svg)](https://codecov.io/gh/s12v/exec-with-secrets)

# Pass secrets from AWS KMS/SSM/Secrets Manager or Azure Key Vault into your app environment

`exec-with-secrets` it supports the following services as secrets providers:
 - [AWS Key Management (KMS)](https://aws.amazon.com/kms/)
 - [AWS Systems Manager Parameter Store (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
 - [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
 - [Azure Key Vault](https://azure.microsoft.com/en-in/services/key-vault/)

This small utility looks for prefixed variables in environment and replaces them with the secret value:
 - `{aws-kms}AQICAHjA3mwbmf...` - decrypts the value using AWS KMS
 - `{aws-ssm}/app/staging/param` - loads parameter `/app/staging/param` from AWS Systems Manager Parameter Store
 - `{aws-sm}/app/staging/param` - loads secret `/app/staging/param` from AWS Secrets Manager
 - `{aws-sm}/app/staging/param{prop1}` - loads secret `/app/staging/param` from AWS Secrets Manager and takes `prop1` property
 - `{az-kv}vault/name` - loads secret `name` from Azure Key Vault `vault`
 
Then it runs `exec` system call. **The secrets are only available to your application and not accessible with `docker inspect`**

Basic example:
```
SECRET="{aws-ssm}/my/secret" exec-with-secrets myapp # $SECRET is plaintext in myapp environment
```

Access:
 - The default credentials chain is used for AWS access
 - Azure authorizer from environment variables/MSI
 - Azure authorizer from configuration file, if the file is set using `AZURE_AUTH_LOCATION` variable

## Examples

### Wrap executable

```
# Download the latest binary
curl -L https://github.com/s12v/exec-with-secrets/releases/download/v0.4.0/exec-with-secrets-darwin-amd64 -o exec-with-secrets
chmod +x ./exec-with-secrets

# Wrap /bin/sh
PARAM="{aws-kms}c2VjcmV0" ./exec-with-secrets /bin/sh -c 'echo $PARAM'
```

`PARAM` will be decrypted and passed to `/bin/sh` via environment.

### Docker example

Build the [example Docker image](Dockerfile):

```
make docker
```

Run:
```
docker run \
    -e PLAINTEXT_PARAM="text" \
    -e KMS_PARAM="{aws-kms}c2VjcmV0" \
    -e SSM_PARAM="{aws-ssm}/myapp/param" \
    exec-with-secrets-example \
    /bin/env
```

`KMS_PARAM` and `SSM_PARAM` will be decrypted and passed to `/bin/env` as environment variables.


## Build

`make` builds Linux and Mac binaries with all providers.

### Choose providers

To chose providers (for example only AWS SSM), run:
```
make TAGS=awsssm
```

## Adding a new provider

See example PR: https://github.com/s12v/exec-with-secrets/pull/1
