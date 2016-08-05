package spec

// PKIConfig is used to configure the setup of a PKI backend done by the
// PKIController.
type PKIConfig struct {
	// AllowedDomains represents a comma separate list of valid domain names the
	// generated certificate authority is valid for.
	AllowedDomains string `json:"allowed_domains"`

	// ClusterID represents the cluster ID a PKI backend setup should be done
	// for. This ID is used to restrict access on Vault related operations for a
	// specific cluster. E.g. the Vault PKI backend will be mounted on a path
	// containing this ID. That way Vault policies will restrict access to this
	// specific path.
	ClusterID string `json:"cluster_id"`

	// CommonName is the common name used to configure the root CA associated
	// with the current PKI backend.
	CommonName string `json:"common_name"`

	// TTL configures the time to live for the root CA being set up. This is a
	// golang time string with the allowed units s, m and h.
	TTL string `json:"ttl"`
}

// PKIController manages the setup of Vault's PKI backends and all other
// required steps necessary to be done.
type PKIController interface {
	// PKI management.

	// SetupPKIBackend sets up a Vault PKI backend according to the given
	// configuration.
	SetupPKIBackend(config PKIConfig) error

	// Path management.

	// CAPath returns the path under which a cluster's certificate authority can
	// be generated. This is very specific to Vault. The path structure is the
	// following. See also https://github.com/hashicorp/vault/blob/6f0f46deb622ba9c7b14b2ec0be24cab3916f3d8/website/source/docs/secrets/pki/index.html.md#pkirootgenerate.
	//
	//     pki-<clusterID>/root/generate/exported
	//
	CAPath(clusterID string) string

	// MountPath returns the path under which a cluster's PKI backend is mounted.
	// This is very specific to Vault. The path structure is the following.
	//
	//     pki-<clusterID>
	//
	MountPath(clusterID string) string

	// RolePath returns the path under which a role is registered. This is very
	// specific to Vault. The path structure is the following. See also https://github.com/hashicorp/vault/blob/6f0f46deb622ba9c7b14b2ec0be24cab3916f3d8/website/source/docs/secrets/pki/index.html.md#pkiroles.
	//
	//     pki-<clusterID>/roles/role-<clusterID>
	//
	RolePath(clusterID string) string
}
