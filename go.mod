module github.com/s12v/exec-with-secrets

require github.com/aws/aws-sdk-go-v2 v0.31.0

require (
	github.com/Azure/azure-sdk-for-go v30.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.13
	github.com/Azure/go-autorest/autorest/azure/auth v0.1.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.2.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v0.31.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v0.31.0
	github.com/aws/aws-sdk-go-v2/service/ssm v0.31.0
)

go 1.13
