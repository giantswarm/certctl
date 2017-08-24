package pki

import (
	"fmt"

	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
)

// ServiceConfig represents the configuration used to create a new PKI controller.
type ServiceConfig struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultServiceConfig provides a default configuration to create a PKI controller.
func DefaultServiceConfig() ServiceConfig {
	newClientConfig := vaultclient.DefaultConfig()
	newClientConfig.Address = "http://127.0.0.1:8200"
	newVaultClient, err := vaultclient.NewClient(newClientConfig)
	if err != nil {
		panic(err)
	}

	newConfig := ServiceConfig{
		// Dependencies.
		VaultClient: newVaultClient,
	}

	return newConfig
}

// NewService creates a new configured PKI controller.
func NewService(config ServiceConfig) (Service, error) {
	// Dependencies.
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "Vault client must not be empty")
	}

	newService := &service{
		ServiceConfig: config,
	}

	return newService, nil
}

type service struct {
	ServiceConfig
}

// PKI management.

func (s *service) Delete(clusterID string) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Unmount the PKI backend, if it exists.
	mounted, err := s.IsMounted(clusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if mounted {
		err = sysBackend.Unmount(s.MountPKIPath(clusterID))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *service) IsCAGenerated(clusterID string) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := s.VaultClient.Logical()

	// Check if a root CA for the given cluster ID exists.
	secret, err := logicalBackend.Read(s.ReadCAPath(clusterID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	// If the secret is nil, the CA has not been generated.
	if secret == nil {
		return false, nil
	}

	if certificate, ok := secret.Data["certificate"]; ok && certificate == "" {
		return false, nil
	}
	if err, ok := secret.Data["error"]; ok && err != "" {
		return false, nil
	}

	return true, nil
}

func (s *service) IsMounted(clusterID string) (bool, error) {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Check if a PKI for the given cluster ID exists.
	mounts, err := sysBackend.ListMounts()
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}
	mountOutput, ok := mounts[s.ListMountsPath(clusterID)+"/"]
	if !ok || mountOutput.Type != "pki" {
		return false, nil
	}

	return true, nil
}

func (s *service) Create(config CreateConfig) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Mount a new PKI backend for the cluster, if it does not already exist.
	mounted, err := s.IsMounted(config.ClusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if !mounted {
		newMountConfig := &vaultclient.MountInput{
			Type:        "pki",
			Description: fmt.Sprintf("PKI backend for cluster ID '%s'", config.ClusterID),
			Config: vaultclient.MountConfigInput{
				MaxLeaseTTL: config.TTL,
			},
		}
		err = sysBackend.Mount(s.MountPKIPath(config.ClusterID), newMountConfig)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's root CA and role.
	logicalBackend := s.VaultClient.Logical()

	// Generate a certificate authority for the PKI backend, if it does not
	// already exist.
	generated, err := s.IsCAGenerated(config.ClusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if !generated {
		data := map[string]interface{}{
			"ttl":         config.TTL,
			"common_name": config.CommonName,
		}
		_, err = logicalBackend.Write(s.WriteCAPath(config.ClusterID), data)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

// Path management.

func (s *service) ReadCAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/cert/ca", clusterID)
}

func (s *service) MountPKIPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (s *service) ListMountsPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (s *service) WriteCAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/root/generate/internal", clusterID)
}
