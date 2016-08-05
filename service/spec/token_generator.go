package spec

type TokenConfig struct {
	// ClusterID represents the cluster ID a token is requested for. This ID is
	// used to restrict access on Vault related operations for a specific
	// cluster. E.g. the generated token will only be allowed to issue
	// certificates for the Vault PKI backend associated with the given cluster
	// ID.
	ClusterID string `json:"cluster_id"`

	// Policies represents the policy names used to be applied to the generated
	// Vault tokens.
	Policies []string `json:"policies"`

	// TTL configures the time to live for the requested token. This is a
	// golang time string with the allowed units s, m and h.
	TTL string `json:"ttl"`
}

// TokenGenerator creates new Vault policies to restrict access capabilities
// of e.g. Vault tokens.
type TokenGenerator interface {
	// NewPKIIssuePolicy creates a new policy to restrict access to only being
	// able to issue signed certificates on the Vault PKI backend specific to the
	// given cluster ID. Here the given cluster ID is used to create the policy
	// name and the policy specific rules matching certain paths within the Vault
	// file system like path structure. This policy name can be used to e.g.
	// apply it to some Vault token. Returned is the policy name and a string
	// representation of the HCL formatted Vault policy. The latter one might not
	// be of interest at all.
	NewPKIIssuePolicy(clusterID string) (string, string, error)

	// NewPKIIssueToken generates a new Vault token allowed to be used to issue
	// signed certificates with respect to the given configuration.
	NewPKIIssueToken(config spec.TokenConfig) (string, error)
}
