package spec

// PolicyGenerator creates new Vault policies to restrict access capabilities
// of e.g. Vault tokens.
type PolicyGenerator interface {
	// NewPKIIssuePolicy creates a new policy to restrict access to only being
	// able to issue signed certificates on the Vault PKI backend specific to the
	// given cluster ID. Here the given cluster ID is used to create the policy
	// name and the policy specific rules matching certain paths within the Vault
	// file system like path structure. This policy name can be used to e.g.
	// apply it to some Vault token. Returned is a string representation of the
	// HCL formatted Vault policy, which might not be of interest at all.
	NewPKIIssuePolicy(clusterID string) (string, error)
}
