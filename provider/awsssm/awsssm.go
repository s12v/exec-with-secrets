// +build awsssm

package awskms

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/s12v/secure-exec/provider"
	"strings"
)

type SsmProvider struct {
	awsSSmClient *ssm.SSM
}

const prefix = "{aws-ssm}"

var fetch func (awsSsmClient *ssm.SSM, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error)

func init() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load AWS-SDK config, " + err.Error())
	}

	fetch = awsFetch
	provider.Register(&SsmProvider{ssm.New(cfg)})
}

func awsFetch(awsSsmClient *ssm.SSM, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	if resp, err := awsSsmClient.GetParameterRequest(input).Send(); err != nil {
		return nil, errors.New(fmt.Sprintf("SSM error: %v", err))
	} else {
		return resp, nil
	}
}

func (p *SsmProvider) Match(val string) bool {
	return strings.HasPrefix(val, prefix) && len(val) > len(prefix)
}

func (p *SsmProvider) Decode(val string) (string, error) {
	name := val[len(prefix):]
	var withEncryption = true
	input := &ssm.GetParameterInput{Name: &name, WithDecryption: &withEncryption}
	if err := input.Validate(); err != nil {
		return "", err
	}

	if output, err := fetch(p.awsSSmClient, input); err != nil {
		return "", err
	} else {
		return *output.Parameter.Value, nil
	}
}
