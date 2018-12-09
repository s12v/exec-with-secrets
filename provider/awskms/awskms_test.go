// +build awskms

package awskms

import (
	"testing"
)

func TestKms_Match(t *testing.T) {
	kms := KmsProvider{}

	if kms.Match("{aws-kms}something") != true {
		t.Fatal("expected to match")
	}

	if kms.Match("https://example.com") != false {
		t.Fatal("not expected to match")
	}
}
