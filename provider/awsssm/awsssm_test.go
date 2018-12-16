// +build awsssm

package awsssm

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"testing"
)

func TestSsmProvider_Match(t *testing.T) {
	ssmProvider := SsmProvider{}

	if ssmProvider.Match("{aws-ssm}something") != true {
		t.Fatal("expected to match")
	}

	if ssmProvider.Match("https://example.com") != false {
		t.Fatal("not expected to match")
	}
}

func TestSsmProvider_Decode(t *testing.T) {
	ssmProvider := SsmProvider{}

	value := "boom"
	fetch = func(awsSsmClient *ssm.SSM, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
		if *input.Name != "/foo/bar" {
			t.Fatalf("unexpected name %v", input.Name)
		}

		return &ssm.GetParameterOutput{Parameter: &ssm.Parameter{Value: &value}}, nil
	}

	if r, _ := ssmProvider.Decode("{aws-ssm}/foo/bar"); r != "boom" {
		t.Fatalf("unexpected plaintext %v", r)
	}
}

func TestSsmProvider_DecodeInvalidInput(t *testing.T) {
	ssmProvider := SsmProvider{}
	r, err := ssmProvider.Decode("{aws-ssm}")
	if err == nil {
		t.Fatal("expected an error", r)
	}
	if r != "" {
		t.Fatalf("unexpected result: '%v'", r)
	}
}
