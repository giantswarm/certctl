package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certctl/service/cert-signer"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type issueFlags struct {
	VaultAddress string
	VaultToken   string

	ClusterID  string
	CommonName string
	TTL        string
}

var (
	issueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Generate signed certificates for a specific cluster.",
		Run:   issueRun,
	}

	newIssueFlags = &issueFlags{}
)

func init() {
	CLICmd.AddCommand(issueCmd)

	issueCmd.Flags().StringVar(&newIssueFlags.VaultAddress, "vault-address", "http://127.0.0.1:8200", "Address used to connect to Vault.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultToken, "vault-token", "", "Token used to authenticate against Vault.")

	issueCmd.Flags().StringVar(&newIssueFlags.ClusterID, "cluster-id", "", "Cluster ID used to generate a new signed certificate for.")
	issueCmd.Flags().StringVar(&newIssueFlags.CommonName, "common-name", "", "Common name used to generate a new signed certificate for.")
	issueCmd.Flags().StringVar(&newIssueFlags.TTL, "ttl", "720h", "TTL used to generate a new signed certificate for.")
}

func issueValidate(newIssueFlags *issueFlags) error {
	if newIssueFlags.VaultToken == "" {
		return maskAnyf(invalidConfigError, "Vault token must not be empty")
	}
	if newIssueFlags.ClusterID == "" {
		return maskAnyf(invalidConfigError, "cluster ID must not be empty")
	}
	if newIssueFlags.CommonName == "" {
		return maskAnyf(invalidConfigError, "common name must not be empty")
	}

	return nil
}

func issueRun(cmd *cobra.Command, args []string) {
	err := issueValidate(newIssueFlags)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.HTTPClient = &http.Client{}
	newVaultFactoryConfig.Address = newIssueFlags.VaultAddress
	newVaultFactoryConfig.AdminToken = newIssueFlags.VaultToken
	newVaultFactory, err := vaultfactory.New(newVaultFactoryConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client and configure it with the provided admin token
	// through the factory.
	newVaultClient, err := newVaultFactory.NewClient()
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a certificate signer to generate a new signed certificate.
	newCertSignerConfig := certsigner.DefaultConfig()
	newCertSignerConfig.VaultClient = newVaultClient
	newCertSigner, err := certsigner.New(newCertSignerConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Generate a new signed certificate.
	newIssueConfig := spec.IssueConfig{
		ClusterID:  newIssueFlags.ClusterID,
		CommonName: newIssueFlags.CommonName,
		TTL:        newIssueFlags.TTL,
	}
	crt, key, ca, err := newCertSigner.Issue(newIssueConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	fmt.Printf("crt: %#v\n", crt)
	fmt.Printf("key: %#v\n", key)
	fmt.Printf("ca: %#v\n", ca)

	// TODO write certificate data into separate files.
}
