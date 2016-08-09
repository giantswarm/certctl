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

func (tg *tokenGenerator) DeletePKIIssuePolicy(clusterID string) error {
	// Get the system backend for policy operations.
	sysBackend := tg.VaultClient.Sys()

	// Delete the policy by name if it is created.
	created, err := tg.IsPKIIssuePolicyCreated(clusterID)
	if err != nil {
		return maskAny(err)
	}
	if created {
		err := sysBackend.DeletePolicy(tg.PKIIssuePolicyName(clusterID))
		if err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func (tg *tokenGenerator) IsPKIIssuePolicyCreated(clusterID string) (bool, error) {
	// Get the system backend for policy operations.
	sysBackend := tg.VaultClient.Sys()

	// Check if the policy is already there.
	policies, err := sysBackend.ListPolicies()
	if err != nil {
		return false, maskAny(err)
	}
	for _, p := range policies {
		if p == tg.PKIIssuePolicyName(clusterID) {
			return true, nil
		}
	}

	return false, nil
}

func (tg *tokenGenerator) NewPKIIssuePolicy(clusterID string) error {
	// Get the system backend for policy operations.
	sysBackend := tg.VaultClient.Sys()

	// Create policy name and HCL policy rules.
	policyName := tg.PKIIssuePolicyName(clusterID)
	rules, err := execTemplate(pkiIssuePolicyTemplate, pkiIssuePolicyContext{ClusterID: clusterID})
	if err != nil {
		return maskAny(err)
	}

	// Actually create the policy within Vault.
	err = sysBackend.PutPolicy(policyName, rules)
	if err != nil {
		return maskAny(err)
	}

	return nil
}

func (tg *tokenGenerator) NewPKIIssueTokens(config spec.TokenConfig) ([]string, error) {
	// In case there does no policy exist that allows to issue certificates on a
	// PKI backend, create one.
	created, err := tg.IsPKIIssuePolicyCreated(config.ClusterID)
	if err != nil {
		return nil, maskAny(err)
	}
	if !created {
		err := tg.NewPKIIssuePolicy(config.ClusterID)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	// Get the token auth backend to create new tokens.
	tokenAuth := tg.VaultClient.Auth().Token()

	// Create the requested amount of tokens.
	var tokens []string
	for i := 0; i < config.Num; i++ {
		tokenID := uuid.New()
		tokens = append(tokens, tokenID)
		newCreateRequest := &vaultclient.TokenCreateRequest{
			ID: tokenID,
			Metadata: map[string]string{
				"cluster-id": config.ClusterID,
			},
			NoParent: true,
			Policies: []string{tg.PKIIssuePolicyName(config.ClusterID)},
			TTL:      config.TTL,
		}
		_, err := tokenAuth.Create(newCreateRequest)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	return tokens, nil
}

func (tg *tokenGenerator) PKIIssuePolicyName(clusterID string) string {
	return fmt.Sprintf("pki-issue-policy-%s", clusterID)
}
