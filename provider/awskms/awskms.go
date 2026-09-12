//go:build awskms
// +build awskms

package awskms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

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
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load AWS SDK config, " + err.Error())
	}

	decrypt = awsDecrypt
	provider.Register(&KmsProvider{kms.NewFromConfig(cfg)})
}

func awsDecrypt(awsKmsClient *kms.Client, input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	ctx := context.Background()
	if resp, err := awsKmsClient.Decrypt(ctx, input); err != nil {
		return nil, fmt.Errorf("KMS error: %w", err)
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

	if len(blob) == 0 {
		return "", errors.New("empty ciphertext")
	}

	input := &kms.DecryptInput{CiphertextBlob: blob}
	if output, err := decrypt(p.awsKmsClient, input); err != nil {
		return "", err
	} else {
		return string(output.Plaintext), nil
	}
}
