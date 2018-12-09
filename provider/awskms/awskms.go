// +build awskms

package awskms

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/s12v/secure-exec/provider"
	"strings"
)

const prefix = "{aws-kms}"

type KmsProvider struct {
	kmsClient awsKms
}

type awsKms interface {
	DecryptRequest(input *kms.DecryptInput) kms.DecryptRequest
}

func init() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load AWS-SDK config, " + err.Error())
	}

	provider.Register(&KmsProvider{kms.New(cfg)})
}

func (p *KmsProvider) Match(val string) bool {
	return strings.HasPrefix(val, prefix) && len(val) > len(prefix)
}

func (p *KmsProvider) Decode(val string) (string, error) {
	blob, err := base64.StdEncoding.DecodeString(val[len(prefix):])
	if err != nil {
		return val, err
	}

	input := &kms.DecryptInput{CiphertextBlob: blob}
	if err = input.Validate(); err != nil {
		return val, nil
	}

	req := p.kmsClient.DecryptRequest(input)
	resp, err := req.Send()
	if err != nil {
		return val, errors.New(fmt.Sprintf("KMS error: %v", err))
	}

	return string(resp.Plaintext), nil
}
