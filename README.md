[![CI](https://github.com/s12v/exec-with-secrets/actions/workflows/ci.yml/badge.svg)](https://github.com/s12v/exec-with-secrets/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/s12v/exec-with-secrets/branch/master/graph/badge.svg)](https://codecov.io/gh/s12v/exec-with-secrets)

# Inject secrets from AWS KMS/SSM/Secrets Manager into your app environment

`exec-with-secrets` supports the following services as secrets providers:
 - [AWS Key Management (KMS)](https://aws.amazon.com/kms/)
 - [AWS Systems Manager Parameter Store (SSM)](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
 - [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)

This utility looks for prefixed variables in environment and replaces them with secret values:
 - `{aws-kms}AQICAHjA3mwbmf...` - decrypts the value using AWS KMS
 - `{aws-ssm}/app/param` - loads parameter `/app/param` from AWS Systems Manager Parameter Store
 - `{aws-sm}/app/param` - loads secret `/app/param` from AWS Secrets Manager
 - `{aws-sm}/app/param[prop1]` - loads secret `/app/param` from AWS Secrets Manager and takes `prop1` property

After decrypting secrets it runs [`exec`](https://en.wikipedia.org/wiki/Exec_(system_call)) system call, replacing itself with your app.
The app can simply access decrypted secrets in the environment.

Basic example:
```
SECRET="{aws-ssm}/my/secret" exec-with-secrets myapp # SECRET value is in myapp environment
```

## Docker example

Build the example Docker image:

```
make docker
```

Run:
```
docker run -e PARAM="text" -e KMS_PARAM="{aws-kms}c2VjcmV0" exec-with-secrets-example sh -c 'echo "$KMS_PARAM"'
```

You need to put a real KMS-encrypted value and pass AWS credentials to the container.

 - `KMS_PARAM` will be decrypted and passed to `sh` as an environment variable
 - `PARAM` will be passed without modifications

You can adapt [Dockerfile](Dockerfile) for your use-case. Use `exec-with-secrets` just like the regular `exec`. For example, run a Java application with:
```
CMD exec-with-secrets java -jar myapp.jar
```
**Note that the decrypted secrets are only visible to your application. `docker inspect` will show encrypted values**

## Secret provider access

Your container should have appropriate permissions to the secrets provider.

The default AWS credentials chain is used.

## Build

`make` builds Linux and macOS binaries (amd64 and arm64) with all providers.

### Choose providers

To chose providers (for example only AWS SSM), run:
```
make TAGS=awsssm
```

## Adding a new provider

See example PR: https://github.com/s12v/exec-with-secrets/pull/1
