package vaultfactory

import (
	"net/http"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certificate-sidekick/service/spec"
)

// Config represents the configuration used to create a new Vault factory.
type Config struct {
	// Dependencies.
	HTTPClient *http.Client

	// Settings.
	Address    string
	AdminToken string
}

// DefaultConfig provides a default configuration to create a Vault factory.
func DefaultConfig() Config {
	newConfig := Config{
		// Dependencies.
		HTTPClient: http.DefaultClient,

		// Settings.
		Address:    "http://127.0.0.1:8200",
		AdminToken: "admin-token",
	}

	return newConfig
}

// New creates a new configured Vault factory.
func New(config Config) (spec.VaultFactory, error) {
	newSecret := &vault{
		Config: config,
	}

	// Dependencies.
	if newSecret.Address == "" {
		return nil, maskAnyf(invalidConfigError, "Vault address must not be empty")
	}
	// Settings.
	if newSecret.HTTPClient == nil {
		return nil, maskAnyf(invalidConfigError, "HTTP client must not be empty")
	}
	if newSecret.AdminToken == "" {
		return nil, maskAnyf(invalidConfigError, "Vault admin token must not be empty")
	}

	return newSecret, nil
}

type vault struct {
	Config
}

func (v *vault) NewClient() (*vaultclient.Client, error) {
	newClientConfig := vaultclient.DefaultConfig()
	newClientConfig.Address = v.Address
	newClientConfig.HttpClient = v.HTTPClient
	newVaultClient, err := vaultclient.NewClient(newClientConfig)
	if err != nil {
		return nil, maskAny(err)
	}
	newVaultClient.SetToken(v.AdminToken)

	return newVaultClient, nil
}
