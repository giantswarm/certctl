package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certctl/service/pki-controller"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certctl/service/token-generator"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type setupFlags struct {
	// Vault
	VaultAddress string
	VaultToken   string

	// Cluster
	ClusterID string

	// PKI
	AllowedDomains string
	CommonName     string
	CATTL          string

	// Token
	NumTokens int
	TokenTTL  string
}

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup a Vault PKI backend including all necessary requirements.",
		Run:   setupRun,
	}

	newSetupFlags = &setupFlags{}
)

func init() {
	CLICmd.AddCommand(setupCmd)

	setupCmd.Flags().StringVar(&newSetupFlags.VaultAddress, "vault-address", "http://127.0.0.1:8200", "Address used to connect to Vault.")
	setupCmd.Flags().StringVar(&newSetupFlags.VaultToken, "vault-token", "", "Token used to authenticate against Vault.")

	setupCmd.Flags().StringVar(&newSetupFlags.ClusterID, "cluster-id", "", "Cluster ID used to generate a new root CA for.")

	setupCmd.Flags().StringVar(&newSetupFlags.AllowedDomains, "allowed-domains", "", "Comma separated domains allowed to authenticate against the cluster's root CA.")
	setupCmd.Flags().StringVar(&newSetupFlags.CommonName, "common-name", "", "Common name used to generate a new root CA for.")
	setupCmd.Flags().StringVar(&newSetupFlags.CATTL, "ca-ttl", "720h", "TTL used to generate a new root CA.")

	setupCmd.Flags().IntVar(&newSetupFlags.NumTokens, "num-tokens", 1, "Number of tokens to generate.")
	setupCmd.Flags().StringVar(&newSetupFlags.TokenTTL, "token-ttl", "720h", "TTL used to generate new tokens.")
}

func setupValidate(newSetupFlags *setupFlags) error {
	if newSetupFlags.VaultToken == "" {
		return maskAnyf(invalidConfigError, "Vault token must not be empty")
	}
	if newSetupFlags.AllowedDomains == "" {
		return maskAnyf(invalidConfigError, "allowed domains must not be empty")
	}
	if newSetupFlags.ClusterID == "" {
		return maskAnyf(invalidConfigError, "cluster ID must not be empty")
	}
	if newSetupFlags.CommonName == "" {
		return maskAnyf(invalidConfigError, "common name must not be empty")
	}

	return nil
}

func setupRun(cmd *cobra.Command, args []string) {
	err := setupValidate(newSetupFlags)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.HTTPClient = &http.Client{}
	newVaultFactoryConfig.Address = newSetupFlags.VaultAddress
	newVaultFactoryConfig.AdminToken = newSetupFlags.VaultToken
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

	// Create a PKI controller to setup the cluster's PKI backend uncluding it's
	// root CA and role.
	newPKIControllerConfig := pkicontroller.DefaultConfig()
	newPKIControllerConfig.VaultClient = newVaultClient
	newPKIController, err := pkicontroller.New(newPKIControllerConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a token generator to create new tokens for the current cluster.
	newTokenGeneratorConfig := tokengenerator.DefaultConfig()
	newTokenGeneratorConfig.VaultClient = newVaultClient
	newTokenGenerator, err := tokengenerator.New(newTokenGeneratorConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Setup PKI backend for cluster.
	newPKIConfig := spec.PKIConfig{
		AllowedDomains: newSetupFlags.AllowedDomains,
		ClusterID:      newSetupFlags.ClusterID,
		CommonName:     newSetupFlags.CommonName,
		TTL:            newSetupFlags.CATTL,
	}
	err = newPKIController.SetupPKIBackend(newPKIConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Generate tokens for the cluster VMs.
	newTokenConfig := spec.TokenConfig{
		ClusterID: newSetupFlags.ClusterID,
		Num:       newSetupFlags.NumTokens,
		TTL:       newSetupFlags.TokenTTL,
	}
	tokens, err := newTokenGenerator.NewPKIIssueTokens(newTokenConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	fmt.Printf("%#v\n", tokens)
}
