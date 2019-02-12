package token

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	"github.com/giantswarm/go-uuid/uuid"
	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
)

const (
	systemMastersOrganizations = "system:masters"
)

// ServiceConfig represents the configuration used to create a new service.
type ServiceConfig struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultServiceConfig provides a default configuration to create a service.
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

// NewService creates a new configured service.
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

func (s *service) Create(config CreateConfig) ([]string, error) {
	// In case there does no policy exist that allows to issue certificates on a
	// PKI backend, create one.
	policyCreated, err := s.IsPolicyCreated(config.ClusterID)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if !policyCreated {
		err := s.CreatePolicy(config.ClusterID)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// In case there is no policy exist that allows to issue certificates
	// with organization on a PKI backend, create one.
	orgPolicyCreated, err := s.IsOrgPolicyCreated(config.ClusterID)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if !orgPolicyCreated {
		err := s.CreateOrgPolicy(config.ClusterID)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// Get the token auth backend to create new tokens.
	tokenAuth := s.VaultClient.Auth().Token()

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
			Policies: []string{s.PolicyName(config.ClusterID)},
			TTL:      config.TTL,
		}
		_, err := tokenAuth.Create(newCreateRequest)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return tokens, nil
}

func (s *service) CreateOrgPolicy(clusterID string) error {
	// Get the system backend for policy operations.
	sysBackend := s.VaultClient.Sys()

	// Create organization policy name and HCL policy rules.
	orgPolicyName := s.OrgPolicyName(clusterID)
	organizationsRoleHash := computeRoleHash(systemMastersOrganizations)
	rules, err := execTemplate(pkiIssueOrgPolicyTemplate, pkiIssueOrgPolicyContext{ClusterID: clusterID, OrganizationsRoleHash: organizationsRoleHash})
	if err != nil {
		return microerror.Mask(err)
	}

	// Actually create the policy within Vault.
	err = sysBackend.PutPolicy(orgPolicyName, rules)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (s *service) CreatePolicy(clusterID string) error {
	// Get the system backend for policy operations.
	sysBackend := s.VaultClient.Sys()

	// Create policy name and HCL policy rules.
	policyName := s.PolicyName(clusterID)
	rules, err := execTemplate(pkiIssuePolicyTemplate, pkiIssuePolicyContext{ClusterID: clusterID})
	if err != nil {
		return microerror.Mask(err)
	}

	// Actually create the policy within Vault.
	err = sysBackend.PutPolicy(policyName, rules)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (s *service) DeleteOrgPolicy(clusterID string) error {
	// Get the system backend for policy operations.
	sysBackend := s.VaultClient.Sys()

	// Delete the policy by name if it is created.
	created, err := s.IsPolicyCreated(clusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if created {
		err := sysBackend.DeletePolicy(s.OrgPolicyName(clusterID))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *service) DeletePolicy(clusterID string) error {
	// Get the system backend for policy operations.
	sysBackend := s.VaultClient.Sys()

	// Delete the policy by name if it is created.
	created, err := s.IsPolicyCreated(clusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if created {
		err := sysBackend.DeletePolicy(s.PolicyName(clusterID))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (s *service) IsPolicyCreated(clusterID string) (bool, error) {
	// Get the system backend for policy operations.
	sysBackend := s.VaultClient.Sys()

	// Check if the policy is already there.
	policies, err := sysBackend.ListPolicies()
	if err != nil {
		return false, microerror.Mask(err)
	}
	for _, p := range policies {
		if p == s.PolicyName(clusterID) {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) IsOrgPolicyCreated(clusterID string) (bool, error) {
	// Get the system backend for policy operations.
	sysBackend := s.VaultClient.Sys()

	// Check if the policy is already there.
	policies, err := sysBackend.ListPolicies()
	if err != nil {
		return false, microerror.Mask(err)
	}
	for _, p := range policies {
		if p == s.OrgPolicyName(clusterID) {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) OrgPolicyName(clusterID string) string {
	return fmt.Sprintf("pki-issue-policy-%s-org", clusterID)
}

func (s *service) PolicyName(clusterID string) string {
	return fmt.Sprintf("pki-issue-policy-%s", clusterID)
}

// computeRoleHash computes a hash for the role that can issue these organizations.
// Since we want to reuse roles when possible, we should try to make sure that
// the same list of organizations returns the same hash (regardless of the order).
// The reason we don't use just the organizations that the user provided is because
// that could potentially be a very long list, or otherwise contain characters
// that are not allowed in URLs.
func computeRoleHash(organizations string) string {
	// Sort organizations alphabetically
	organizationsSlice := strings.Split(organizations, ",")
	sort.Strings(organizationsSlice)
	organizations = strings.Join(organizationsSlice, ",")

	h := sha1.New()
	h.Write([]byte(organizations))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
