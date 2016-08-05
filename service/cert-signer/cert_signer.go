package certsigner

import (
	"fmt"
	"net/http"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/spec"
)

// Config represents the configuration used to create a new certificate signer.
type Config struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultConfig provides a default configuration to create a certificate
// signer.
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

// New creates a new configured certificate signer.
func New(config Config) (spec.PolicyGenerator, error) {
	newCertSigner := &certSigner{
		Config: config,
	}

	// Dependencies.
	if newCertSigner.VaultClient == nil {
		return nil, maskAnyf(invalidConfigError, "Vault client must not be empty")
	}

	return newCertSigner, nil
}

type certSigner struct {
	Config
}

func (cs *certSigner) Issue(config spec.IssueConfig) error {
	// Create a client for issuing a new signed certificate.
	logicalStore := newVaultClient.Logical()

	// Generate a certificate for the PKI backend signed by the certificate
	// authority associated with the configured cluster ID.
	data := map[string]interface{}{
		"ttl":         config.TTL,
		"common_name": config.CommonName,
	}
	secret, err := logicalStore.Write(cs.SignedPath(config.ClusterID), data)
	if err != nil {
		return maskAny(err)
	}

	return nil
}

func (cs *certSigner) SignedPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/issue/role-%s", clusterID, clusterID)
}
