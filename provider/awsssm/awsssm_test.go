// +build awsssm

package awskms

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"testing"
)

func TestSsmProvider_Match(t *testing.T) {
	kmsProvider := SsmProvider{}

	if kmsProvider.Match("{aws-ssm}something") != true {
		t.Fatal("expected to match")
	}

	if kmsProvider.Match("https://example.com") != false {
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
