// +build awsssm

package awsssm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
	fetch = awsFetch

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		fmt.Println("error obtaining AWS credentials:", err)
		os.Exit(1)
	}

	opts := ssm.Options{
		Credentials: cfg.Credentials,
	}

	ssmClient := ssm.New(opts)

	provider.Register(&SsmProvider{ssmClient})
}

func awsFetch(awsSsmClient *ssm.Client, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := awsSsmClient.GetParameter(ctx, input)
	if err != nil {
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
	input := &ssm.GetParameterInput{Name: &name, WithDecryption: true}

	if output, err := fetch(p.awsSsmClient, input); err != nil {
		return "", fmt.Errorf("failed to fetch secret %s: %s", name, err)
	} else {
		return *output.Parameter.Value, nil
	}
}
