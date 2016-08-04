package policygenerator

import (
	"net/http"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/spec"
)

// Config represents the configuration used to create a new policy generator.
type Config struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultConfig provides a default configuration to create a policy generator.
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

// New creates a new configured policy generator.
func New(config Config) (spec.PolicyGenerator, error) {
	newPolicyGenerator := &policyGenerator{
		Config: config,
	}

	// Dependencies.
	if newPolicyGenerator.VaultClient == nil {
		return nil, maskAnyf(invalidConfigError, "Vault client must not be empty")
	}

	return newPolicyGenerator, nil
}

type policyGenerator struct {
	Config
}

func (pg *policyGenerator) NewPKIIssuePolicy(clusterID string) (string, error) {
	// create HCL policy string
	rules, err := execTemplate(pkiIssuePolicyTemplate, pkiIssuePolicyContext{ClusterID: clusterID})
	if err != nil {
		return "", maskAny(err)
	}

	sysBackend := pg.VaultClient.Sys()
	err = sysBackend.PutPolicy(clusterID, rules)
	if err != nil {
		return "", maskAny(err)
	}

	return rules, nil
}
