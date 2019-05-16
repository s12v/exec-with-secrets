// +build awskms

package awskms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/s12v/exec-with-secrets/provider"
	"strings"
)

type KmsProvider struct {
	awsKmsClient *kms.KMS
}

const prefix = "{aws-kms}"

var decrypt func(awsKmsClient *kms.KMS, input *kms.DecryptInput) (*kms.DecryptOutput, error)

func init() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load AWS-SDK config, " + err.Error())
	}

	decrypt = awsDecrypt
	provider.Register(&KmsProvider{kms.New(cfg)})
}

func awsDecrypt(awsKmsClient *kms.KMS, input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	ctx := context.Background()
	if resp, err := awsKmsClient.DecryptRequest(input).Send(ctx); err != nil {
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
	if err = input.Validate(); err != nil {
		return "", err
	}

	if output, err := decrypt(p.awsKmsClient, input); err != nil {
		return "", err
	} else {
		return string(output.Plaintext), nil
	}
}
