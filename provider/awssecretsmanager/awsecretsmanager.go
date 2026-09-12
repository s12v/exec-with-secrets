//go:build awssecretsmanager
// +build awssecretsmanager

package awssecretsmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/s12v/exec-with-secrets/provider"
)

type SecretsManagerProvider struct {
	awsClient *secretsmanager.Client
}

const prefix = "{aws-sm}"

var postfix = regexp.MustCompile(`\[[^]]+\]$`)

var fetch func(
	awsClient *secretsmanager.Client,
	input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load AWS SDK config, " + err.Error())
	}

	fetch = awsFetch
	provider.Register(&SecretsManagerProvider{secretsmanager.NewFromConfig(cfg)})
}

func awsFetch(
	awsClient *secretsmanager.Client,
	input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	ctx := context.Background()
	if resp, err := awsClient.GetSecretValue(ctx, input); err != nil {
		return nil, fmt.Errorf("AWS SecretsManager error: %w", err)
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
		return p.decodeJson(name, strings.Trim(property, "[]"))
	}
	return p.fetchString(name)
}

func (p *SecretsManagerProvider) decodeJson(val string, property string) (string, error) {
	name := val[:len(val)-len(property)-2]
	jsobj, err := p.fetchString(name)
	if err != nil {
		return "", err
	}

	properties, err := unmarshal(jsobj)
	if err != nil {
		return "", fmt.Errorf("secret '%v' is not a flat JSON object: %w", name, err)
	}

	value, ok := properties[property]
	if !ok {
		return "", fmt.Errorf("property '%v' does not exist", property)
	}
	return value, nil
}

func (p *SecretsManagerProvider) fetchString(name string) (string, error) {
	if name == "" {
		return "", errors.New("missing secret id")
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
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
