// +build awssecretsmanager

package awssecretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/johnrichardrinehart/exec-with-secrets/provider"
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
	fetch = awsFetch // set global

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("default"),
	)
	if err != nil {
		fmt.Println("error obtaining AWS credentials:", err)
		os.Exit(1)
	}

	opts := secretsmanager.Options{
		Region:      "us-east-1",
		Credentials: cfg.Credentials,
	}
	smClient := secretsmanager.New(opts)

	provider.Register(&SecretsManagerProvider{smClient}) // register externally
}

func awsFetch(
	awsClient *secretsmanager.Client,
	input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := awsClient.GetSecretValue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("AWS SecretsManager error: %v", err)
	}
	return resp, nil
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

	properties, _ := unmarshal(jsobj)
	value, ok := properties[property]
	if !ok {
		return "", fmt.Errorf("property '%v' does not exist", property)
	}
	return value, nil
}

func (p *SecretsManagerProvider) fetchString(name string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	if output, err := fetch(p.awsClient, input); err != nil {
		return "", fmt.Errorf("failed to fetch secret %s: %s", name, err)
	} else {
		return *output.SecretString, nil
	}
}

func unmarshal(val string) (map[string]string, error) {
	var omap map[string]string
	err := json.Unmarshal([]byte(val), &omap)
	return omap, err
}
