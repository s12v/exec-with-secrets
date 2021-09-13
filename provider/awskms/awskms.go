// +build awskms

package awskms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/s12v/exec-with-secrets/provider"
)

type KmsProvider struct {
	awsKmsClient *kms.Client
}

const prefix = "{aws-kms}"

var decrypt func(awsKmsClient *kms.Client, input *kms.DecryptInput) (*kms.DecryptOutput, error)

func init() {
	decrypt = awsDecrypt

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		fmt.Println("error obtaining AWS credentials:", err)
		os.Exit(1)
	}

	opts := kms.Options{
		Region:      "us-east-1",
		Credentials: cfg.Credentials,
	}

	kmsClient := kms.New(opts)

	provider.Register(&KmsProvider{kmsClient})
}

func awsDecrypt(awsKmsClient *kms.Client, input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := awsKmsClient.Decrypt(ctx, input)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("KMS error: %v", err))
	} else {
		return resp, nil
	}
}

func (p *KmsProvider) Match(val string) bool {
	return strings.HasPrefix(val, prefix) && len(val) > len(prefix)
}

func (p *KmsProvider) Decode(val string) (string, error) {
	blob, err := base64.StdEncoding.DecodeString(val[len(prefix):])
	if err != nil {
		return "", err
	}

	input := &kms.DecryptInput{CiphertextBlob: blob}

	if output, err := decrypt(p.awsKmsClient, input); err != nil {
		return "", fmt.Errorf("failed to decrypt secret %s: %s", val, err)
	} else {
		return string(output.Plaintext), nil
	}
}
