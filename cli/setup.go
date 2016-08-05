package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certctl/service/token-generator"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type setupFlags struct {
	VaultAddress string
	VaultToken   string
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
}

func setupRun(cmd *cobra.Command, args []string) {
	if newSetupFlags.VaultToken == "" {
		log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "Vault token must not be empty"))
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

	// Create a token generator to create new tokens for the current cluster.
	newTokenGeneratorConfig := tokengenerator.DefaultConfig()
	newTokenGeneratorConfig.VaultClient = newVaultClient
	newTokenGenerator, err := tokengenerator.New(newTokenGeneratorConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	fmt.Printf("%#v\n", newTokenGenerator)
	// TODO
}
