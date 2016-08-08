package pkicontroller

import (
	"fmt"
	"net/http"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/spec"
)

// Config represents the configuration used to create a new PKI controller.
type Config struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultConfig provides a default configuration to create a PKI controller.
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

// New creates a new configured PKI controller.
func New(config Config) (spec.PKIController, error) {
	newPKIController := &pkiController{
		Config: config,
	}

	// Dependencies.
	if newPKIController.VaultClient == nil {
		return nil, maskAnyf(invalidConfigError, "Vault client must not be empty")
	}

	return newPKIController, nil
}

type pkiController struct {
	Config
}

// PKI management.

func (pc *pkiController) SetupPKIBackend(config spec.PKIConfig) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := pc.VaultClient.Sys()
	// Mount a new PKI backend for the cluster, if it does not already exist.
	mounts, err := sysBackend.ListMounts()
	if err != nil {
		return maskAny(err)
	}
	mountOutput, ok := mounts[pc.MountPath(config.ClusterID)+"/"]
	if !ok || mountOutput.Type != "pki" {
		newMountConfig := &vaultclient.MountInput{
			Type:        "pki",
			Description: fmt.Sprintf("PKI backend for cluster ID '%s'", config.ClusterID),
			Config: vaultclient.MountConfigInput{
				MaxLeaseTTL: config.TTL,
			},
		}
		err = sysBackend.Mount(pc.MountPath(config.ClusterID), newMountConfig)
		if err != nil {
			return maskAny(err)
		}
	}

	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's root CA and role.
	logicalStore := pc.VaultClient.Logical()
	// Generate a certificate authority for the PKI backend.
	data := map[string]interface{}{
		"ttl":         config.TTL,
		"common_name": config.CommonName,
	}
	_, err = logicalStore.Write(pc.CAPath(config.ClusterID), data)
	if err != nil {
		return maskAny(err)
	}
	// Create role for the mounted PKI backend.
	data = map[string]interface{}{
		"allowed_domains":  config.AllowedDomains,
		"allow_subdomains": "true",
		"ttl":              config.TTL,
	}
	_, err = logicalStore.Write(pc.RolePath(config.ClusterID), data)
	if err != nil {
		return maskAny(err)
	}

	return nil
}

// Path management.

func (pc *pkiController) CAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/root/generate/exported", clusterID)
}

func (pc *pkiController) MountPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (pc *pkiController) RolePath(clusterID string) string {
	return fmt.Sprintf("pki-%s/roles/role-%s", clusterID, clusterID)
}
