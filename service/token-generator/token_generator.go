package tokengenerator

import (
	"fmt"
	"net/http"

	"github.com/giantswarm/go-uuid/uuid"
	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/spec"
)

// Config represents the configuration used to create a new token generator.
type Config struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultConfig provides a default configuration to create a token generator.
func DefaultConfig() Config {
	newClientConfig := vaultclient.DefaultConfig()
	newClientConfig.Address = "http://127.0.0.1:8200"
	newClientConfig.HttpClient = http.DefaultClient
	newVaultClient, err := vaultclient.NewClient(newClientConfig)
	if err != nil {
		panic(err)
	}

	newConfig := Config{
		// Dependencies.
		VaultClient: newVaultClient,
	}

	return newConfig
}

// New creates a new configured token generator.
func New(config Config) (spec.TokenGenerator, error) {
	newTokenGenerator := &tokenGenerator{
		Config: config,
	}

	// Dependencies.
	if newTokenGenerator.VaultClient == nil {
		return nil, maskAnyf(invalidConfigError, "Vault client must not be empty")
	}

	return newTokenGenerator, nil
}

type tokenGenerator struct {
	Config
}

func (tg *tokenGenerator) NewPKIIssuePolicy(clusterID string) (string, string, error) {
	// Create HCL policy rules.
	rules, err := execTemplate(pkiIssuePolicyTemplate, pkiIssuePolicyContext{ClusterID: clusterID})
	if err != nil {
		return "", maskAny(err)
	}

	// Actually create the policy within Vault.
	sysBackend := tg.VaultClient.Sys()
	policyName := uuid.New()
	err = sysBackend.PutPolicy(policyName, rules)
	if err != nil {
		return "", maskAny(err)
	}

	return policyName, rules, nil
}

func (tg *tokenGenerator) NewPKIIssueToken(config spec.TokenConfig) (string, error) {
	tokenAuth := tg.VaultClient.Auth().Token()
	tokenName := uuid.New()
	newCreateRequest := &vaultclient.TokenCreateRequest{
		ID: tokenName,
		Metadata: map[string]string{
			"clusterid": config.ClusterID,
		},
		Policies: config.Policies,
		TTL:      config.TTL,
	}
	_, err := tokenAuth.Create(newCreateRequest)
	if err != nil {
		return "", maskAny(err)
	}

	return tokenName, nil
}
