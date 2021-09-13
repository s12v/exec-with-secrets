// +build azurekeyvault

package azurekeyvault

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/johnrichardrinehart/exec-with-secrets/provider"
	"os"
	"strings"
)

type KvProvider struct {
	client *keyvault.BaseClient
}

const prefix = "{az-kv}"

var getSecret func(basicClient *keyvault.BaseClient, vaultName string, name string) (string, error)

func init() {
	getSecret = azureGetSecret

	var authorizer autorest.Authorizer
	var err error
	if _, exists := os.LookupEnv("AZURE_AUTH_LOCATION"); exists {
		authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	} else {
		authorizer, err = auth.NewAuthorizerFromEnvironment()
	}

	if err != nil {
		panic("unable to create Azure Key Vault authorizer: " + err.Error())
	}

	kvClient := keyvault.New()
	kvClient.Authorizer = authorizer
	provider.Register(&KvProvider{&kvClient})
}

func azureGetSecret(basicClient *keyvault.BaseClient, vaultName string, name string) (string, error) {
	ctx := context.Background()
	secretResp, err := basicClient.GetSecret(ctx, vaultBaseUrl(vaultName), name, "")
	if err != nil {
		return "", err
	}

	return *secretResp.Value, nil
}

func vaultBaseUrl(vaultName string) string {
	return "https://" + vaultName + ".vault.azure.net"
}

func (p *KvProvider) Match(val string) bool {
	return strings.HasPrefix(val, prefix) && len(val) > len(prefix) && strings.Contains(val, "/")
}

func (p *KvProvider) Decode(val string) (string, error) {
	vault, name, err := splitPath(val[len(prefix):])

	if err != nil {
		return "", err
	}

	if secret, err := getSecret(p.client, vault, name); err != nil {
		return "", err
	} else {
		return secret, nil
	}
}

func splitPath(compound string) (vault string, name string, err error) {
	parts := strings.Split(compound, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid vault path")
	}

	return parts[0], parts[1], nil
}
