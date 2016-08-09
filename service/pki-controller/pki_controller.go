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

func (pc *pkiController) DeletePKIBackend(clusterID string) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := pc.VaultClient.Sys()

	// Unmount the PKI backend, if it exists.
	mounted, err := pc.IsPKIMounted(clusterID)
	if err != nil {
		return maskAny(err)
	}
	if mounted {
		err = sysBackend.Unmount(pc.MountPKIPath(clusterID))
		if err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func (pc *pkiController) IsCAGenerated(clusterID string) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := pc.VaultClient.Logical()

	// Check if a root CA for the given cluster ID exists.
	secret, err := logicalBackend.Read(pc.ReadCAPath(clusterID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, maskAny(err)
	}
	if certificate, ok := secret.Data["certificate"]; ok && certificate == "" {
		return false, nil
	}
	if err, ok := secret.Data["error"]; ok && err != "" {
		return false, nil
	}

	return true, nil
}

func (pc *pkiController) IsPKIMounted(clusterID string) (bool, error) {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := pc.VaultClient.Sys()

	// Check if a PKI for the given cluster ID exists.
	mounts, err := sysBackend.ListMounts()
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, maskAny(err)
	}
	mountOutput, ok := mounts[pc.ListMountsPath(clusterID)+"/"]
	if !ok || mountOutput.Type != "pki" {
		return false, nil
	}

	return true, nil
}

func (pc *pkiController) IsRoleCreated(clusterID string) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := pc.VaultClient.Logical()

	// Check if a PKI for the given cluster ID exists.
	secret, err := logicalBackend.List(pc.ListRolesPath(clusterID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, maskAny(err)
	}

	// In case there is not a single role for this PKI backend, secret is nil.
	if secret == nil {
		return false, nil
	}

	// When listing roles a list of role names is returned. Here we iterate over
	// this list and if we find the desired role name, it means the role has
	// already been created.
	if keys, ok := secret.Data["keys"]; ok {
		if list, ok := keys.([]interface{}); ok {
			for _, k := range list {
				if s, ok := k.(string); ok && s == pc.PKIRoleName(clusterID) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (pc *pkiController) PKIRoleName(clusterID string) string {
	return fmt.Sprintf("role-%s", clusterID)
}

func (pc *pkiController) SetupPKIBackend(config spec.PKIConfig) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := pc.VaultClient.Sys()

	// Mount a new PKI backend for the cluster, if it does not already exist.
	mounted, err := pc.IsPKIMounted(config.ClusterID)
	if err != nil {
		return maskAny(err)
	}
	if !mounted {
		newMountConfig := &vaultclient.MountInput{
			Type:        "pki",
			Description: fmt.Sprintf("PKI backend for cluster ID '%s'", config.ClusterID),
			Config: vaultclient.MountConfigInput{
				MaxLeaseTTL: config.TTL,
			},
		}
		err = sysBackend.Mount(pc.MountPKIPath(config.ClusterID), newMountConfig)
		if err != nil {
			return maskAny(err)
		}
	}

	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's root CA and role.
	logicalBackend := pc.VaultClient.Logical()

	// Generate a certificate authority for the PKI backend, if it does not
	// already exist.
	generated, err := pc.IsCAGenerated(config.ClusterID)
	if err != nil {
		return maskAny(err)
	}
	if !generated {
		data := map[string]interface{}{
			"ttl":         config.TTL,
			"common_name": config.CommonName,
		}
		_, err = logicalBackend.Write(pc.WriteCAPath(config.ClusterID), data)
		if err != nil {
			return maskAny(err)
		}
	}

	// Create a role for the mounted PKI backend, if it does not already exist.
	created, err := pc.IsRoleCreated(config.ClusterID)
	if err != nil {
		return maskAny(err)
	}
	if !created {
		data := map[string]interface{}{
			"allowed_domains":  config.AllowedDomains,
			"allow_subdomains": "true",
			"ttl":              config.TTL,
		}
		_, err = logicalBackend.Write(pc.WriteRolePath(config.ClusterID), data)
		if err != nil {
			return maskAny(err)
		}
	}

	return nil
}

// Path management.

func (pc *pkiController) ReadCAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/cert/ca", clusterID)
}

func (pc *pkiController) MountPKIPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (pc *pkiController) ListMountsPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (pc *pkiController) ListRolesPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/roles/", clusterID)
}

func (pc *pkiController) WriteCAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/root/generate/internal", clusterID)
}

func (pc *pkiController) WriteRolePath(clusterID string) string {
	return fmt.Sprintf("pki-%s/roles/%s", clusterID, pc.PKIRoleName(clusterID))
}
