module github.com/s12v/exec-with-secrets

require github.com/aws/aws-sdk-go-v2 v1.16.8

require (
	github.com/Azure/azure-sdk-for-go v30.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.19
	github.com/Azure/go-autorest/autorest/azure/auth v0.1.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.2.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.18.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.15.14
	github.com/aws/aws-sdk-go-v2/service/ssm v1.27.6
	github.com/stretchr/testify v1.4.0 // indirect
)

go 1.13
