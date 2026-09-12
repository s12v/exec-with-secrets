//go:build awsssm
// +build awsssm

package awsssm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/s12v/exec-with-secrets/provider"
)

type SsmProvider struct {
	awsSsmClient *ssm.Client
}

const prefix = "{aws-ssm}"

var fetch func(awsSsmClient *ssm.Client, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load AWS SDK config, " + err.Error())
	}

	fetch = awsFetch
	provider.Register(&SsmProvider{ssm.NewFromConfig(cfg)})
}

func awsFetch(awsSsmClient *ssm.Client, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	ctx := context.Background()
	if resp, err := awsSsmClient.GetParameter(ctx, input); err != nil {
		return nil, fmt.Errorf("SSM error: %w", err)
	} else {
		return resp, nil
	}
}

func (p *SsmProvider) Match(val string) bool {
	return strings.HasPrefix(val, prefix) && len(val) > len(prefix)
}

func (p *SsmProvider) Decode(val string) (string, error) {
	name := val[len(prefix):]
	if name == "" {
		return "", errors.New("missing parameter name")
	}

	input := &ssm.GetParameterInput{Name: aws.String(name), WithDecryption: aws.Bool(true)}
	if output, err := fetch(p.awsSsmClient, input); err != nil {
		return "", err
	} else {
		return *output.Parameter.Value, nil
	}
}
