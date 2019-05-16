// +build awssecretsmanager

package awssecretsmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/s12v/exec-with-secrets/provider"
	"regexp"
	"strings"
)

type SecretsManagerProvider struct {
	awsClient *secretsmanager.SecretsManager
}

const prefix = "{aws-sm}"

var postfix = regexp.MustCompile("{[^{^}]+}$")

var fetch func(
	awsClient *secretsmanager.SecretsManager,
	input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)

func init() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load AWS-SDK config, " + err.Error())
	}

	fetch = awsFetch
	provider.Register(&SecretsManagerProvider{secretsmanager.New(cfg)})
}

func awsFetch(
	awsClient *secretsmanager.SecretsManager,
	input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	ctx := context.Background()
	if resp, err := awsClient.GetSecretValueRequest(input).Send(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("AWS SecretsManager error: %v", err))
	} else {
		return resp, nil
	}
}

func (p *SecretsManagerProvider) Match(val string) bool {
	return strings.HasPrefix(val, prefix) && len(val) > len(prefix)
}

func (p *SecretsManagerProvider) Decode(val string) (string, error) {
	name := val[len(prefix):]
	property := postfix.FindString(name)
	if property != "" {
		return p.decodeJson(name, strings.Trim(property, "{}"))
	}
	return p.fetchString(name)
}

func (p *SecretsManagerProvider) decodeJson(val string, property string) (string, error) {
	name := val[:len(val)-len(property)-2]
	jsobj, err := p.fetchString(name)
	if err != nil {
		return "", err
	}

	properties, _ := unmarshal(jsobj)
	value, ok := properties[property]
	if !ok {
		return "", errors.New(fmt.Sprintf("property '%v' does not exist", property))
	}
	return value, nil
}

func (p *SecretsManagerProvider) fetchString(name string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}
	if err := input.Validate(); err != nil {
		return "", err
	}

	if output, err := fetch(p.awsClient, input); err != nil {
		return "", err
	} else {
		return *output.SecretString, nil
	}
}

func unmarshal(val string) (map[string]string, error) {
	var omap map[string]string
	err := json.Unmarshal([]byte(val), &omap)
	return omap, err
}
