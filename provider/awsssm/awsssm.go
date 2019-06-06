// +build awsssm

package awsssm

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/s12v/exec-with-secrets/provider"
	"strings"
)

type SsmProvider struct {
	awsSsmClient *ssm.Client
}

const prefix = "{aws-ssm}"

var fetch func(awsSsmClient *ssm.Client, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error)

func init() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load AWS-SDK config, " + err.Error())
	}

	fetch = awsFetch
	provider.Register(&SsmProvider{ssm.New(cfg)})
}

func awsFetch(awsSsmClient *ssm.Client, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	ctx := context.Background()
	if resp, err := awsSsmClient.GetParameterRequest(input).Send(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("SSM error: %v", err))
	} else {
		return resp.GetParameterOutput, nil
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

	if output, err := fetch(p.awsSsmClient, input); err != nil {
		return "", err
	} else {
		return *output.Parameter.Value, nil
	}
}
