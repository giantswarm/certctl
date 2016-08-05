package spec

// IssueConfig is used to configure the process of issuing a certificate key
// pair using the CertSigner.
type IssueConfig struct {
	// TODO

	// ClusterID represents the cluster ID a PKI backend setup should be done.
	// This ID is used to restrict access on Vault related operations for a
	// specific cluster. E.g. the Vault PKI backend will be mounted on a path
	// containing this ID. That way Vault policies will restrict access to this
	// specific path.
	ClusterID string `json:"cluster_id"`

	// CommonName is the common name used to configure the root CA associated
	// with the current PKI backend.
	CommonName string `json:"common_name"`

	// TTL configures the time to live for the requested certificate. This is a
	// golang time string with the allowed units s, m and h.
	TTL string `json:"ttl"`
}

// CertSigner manages the process of issuing new certificate key pairs
type CertSigner interface {
	Issue()
}
