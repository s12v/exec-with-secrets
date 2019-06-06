// +build awskms

package awskms

import (
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"testing"
)

func init() {
	decrypt = func(awsKmsClient *kms.Client, input *kms.DecryptInput) (*kms.DecryptOutput, error) {
		return &kms.DecryptOutput{Plaintext: input.CiphertextBlob}, nil
	}
}

func TestKmsProvider_Match(t *testing.T) {
	kmsProvider := KmsProvider{}

	if kmsProvider.Match("{aws-kms}something") != true {
		t.Fatal("expected to match")
	}

	if kmsProvider.Match("https://example.com") != false {
		t.Fatal("not expected to match")
	}
}

func TestKmsProvider_Decode(t *testing.T) {
	kmsProvider := KmsProvider{}

	if r, _ := kmsProvider.Decode("{aws-kms}Ym9vbQ=="); r != "boom" {
		t.Fatalf("unexpected plaintext %v", r)
	}
}

func TestKmsProvider_DecodeInvalid(t *testing.T) {
	kmsProvider := KmsProvider{}
	if _, err := kmsProvider.Decode("{aws-kms}"); err == nil {
		t.Fatal("expected an error")
	}
}
