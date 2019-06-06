// +build awssecretsmanager

package awssecretsmanager

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"testing"
)

func TestSecretsManagerProvider_Match(t *testing.T) {
	provider := SecretsManagerProvider{}

	if provider.Match("{aws-sm}something") != true {
		t.Fatal("expected to match")
	}

	if provider.Match("https://example.com") != false {
		t.Fatal("not expected to match")
	}
}

func TestSecretsManagerProvider_Decode(t *testing.T) {
	provider := SecretsManagerProvider{}

	value := "boom"
	fetch = func(
		awsClient *secretsmanager.Client,
		input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
		if *input.SecretId != "/foo/bar" {
			t.Fatalf("unexpected SecretId %v", input.SecretId)
		}

		return &secretsmanager.GetSecretValueOutput{SecretString: &value}, nil
	}

	if r, _ := provider.Decode("{aws-sm}/foo/bar"); r != "boom" {
		t.Fatalf("unexpected value %v", r)
	}
}

func TestSecretsManagerProvider_DecodeJson(t *testing.T) {
	provider := SecretsManagerProvider{}

	value := `{"prop1": "aaa", "prop2": "bbb"}`
	fetch = func(
		awsClient *secretsmanager.Client,
		input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
		if *input.SecretId != "/foo/bar" {
			t.Fatalf("unexpected SecretId %v", *input.SecretId)
		}

		return &secretsmanager.GetSecretValueOutput{SecretString: &value}, nil
	}

	if r, _ := provider.Decode("{aws-sm}/foo/bar{prop2}"); r != "bbb" {
		t.Fatalf("unexpected value %v", r)
	}
}

func TestSecretsManagerProvider_DecodeJson_MissingProperty(t *testing.T) {
	provider := SecretsManagerProvider{}

	value := `{"prop1": "foo", "prop2": "bar"}`
	fetch = func(
		awsClient *secretsmanager.Client,
		input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
		if *input.SecretId != "/foo/bar" {
			t.Fatalf("unexpected SecretId %v", *input.SecretId)
		}

		return &secretsmanager.GetSecretValueOutput{SecretString: &value}, nil
	}

	if _, err := provider.Decode("{aws-sm}/foo/bar{prop3}"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestSecretsManagerProvider_Decode_FetchError(t *testing.T) {
	provider := SecretsManagerProvider{}

	fetch = func(
		awsClient *secretsmanager.Client,
		input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {

		return nil, errors.New("test error")
	}

	if _, err := provider.Decode("{aws-sm}/foo/bar"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestSecretsManagerProvider_DecodeJson_FetchError(t *testing.T) {
	provider := SecretsManagerProvider{}

	fetch = func(
		awsClient *secretsmanager.Client,
		input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {

		return nil, errors.New("test error")
	}

	if _, err := provider.Decode("{aws-sm}/foo/bar{prop1}"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestSecretsManagerProvider_Decode_InvalidInput(t *testing.T) {
	provider := SecretsManagerProvider{}
	r, err := provider.Decode("{aws-sm}")
	if err == nil {
		t.Fatal("expected an error", r)
	}
	if r != "" {
		t.Fatalf("unexpected result: '%v'", r)
	}
}

func Test_Unmarshal(t *testing.T) {
	jsonobj, err := unmarshal(`{"prop1": "foo", "prop2": "bar"}`)

	if err != nil {
		t.Fatal("expected error: ", err)
	}

	if jsonobj["prop1"] != "foo" {
		t.Fatalf("unexpected value '%v'", jsonobj["prop1"])
	}

	if jsonobj["prop2"] != "bar" {
		t.Fatalf("unexpected value '%v'", jsonobj["prop2"])
	}
}
