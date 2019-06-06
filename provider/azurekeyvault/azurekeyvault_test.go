// +build azurekeyvault

package azurekeyvault

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"testing"
)

func Test_SplitPath(t *testing.T) {
	vault, name, err := splitPath("test-vault/test-secret")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if vault != "test-vault" {
		t.Fatal("unexpected vault name:", vault)
	}
	if name != "test-secret" {
		t.Fatal("unexpected secret name:", name)
	}
}

func Test_SplitPathBadInput(t *testing.T) {
	if _, _, err := splitPath("test-secret"); err == nil {
		t.Fatal("expected an error")
	}
}

func Test_Match(t *testing.T) {
	kvProvider := KvProvider{}
	if kvProvider.Match("{az-kv}my-vaulet/my-secret") != true {
		t.Fatal("expected to match")
	}

	if kvProvider.Match("{az-kv}my-secret") != false {
		t.Fatal("not expected to match")
	}

	if kvProvider.Match("my-secret") != false {
		t.Fatal("not expected to match")
	}
}

func Test_Decode(t *testing.T) {
	kvProvider := KvProvider{}
	getSecret = func(basicClient *keyvault.BaseClient, vaultName string, name string) (string, error) {
		if vaultName != "vault1" {
			t.Fatal("unexpected vault name:", vaultName)
		}

		return "secretvalue", nil
	}

	if r, err := kvProvider.Decode("{az-kv}vault1/bar"); r != "secretvalue" {
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		t.Fatal("unexpected plaintext:", r)
	}
}

func Test_DecodeInvalidInput(t *testing.T) {
	kvProvider := KvProvider{}
	r, err := kvProvider.Decode("{az-kv}some-name")
	if err == nil {
		t.Fatal("expected an error", r)
	}
	if r != "" {
		t.Fatal("unexpected result:", r)
	}
}

func Test_VaultBaseName(t *testing.T) {
	url := vaultBaseUrl("test1")
	if url != "https://test1.vault.azure.net" {
		t.Fatal("unexpected vault url:", url)
	}
}
