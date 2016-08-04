package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certctl/service/policy-generator"
	"github.com/giantswarm/certctl/service/vault-factory"
)

var (
	generateSignedCmd = &cobra.Command{
		Use:   "signed",
		Short: "Generate signed certificates.",
		Run:   generateSignedRun,
	}

	generateSignedVaultAddress string
	generateSignedVaultToken   string
)

func init() {
	generateCmd.AddCommand(generateSignedCmd)

	generateSignedCmd.Flags().StringVar(&generateSignedVaultAddress, "vault-address", "http://127.0.0.1:8200", "Address used to connect to Vault.")
	generateSignedCmd.Flags().StringVar(&generateSignedVaultToken, "vault-token", "", "Token used to authenticate against Vault.")
}

func generateSignedRun(cmd *cobra.Command, args []string) {
	if generateSignedVaultToken == "" {
		log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "Vault token must not be empty"))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.HTTPClient = &http.Client{}
	newVaultFactoryConfig.Address = generateSignedVaultAddress
	newVaultFactoryConfig.AdminToken = generateSignedVaultToken
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

	newPolicyGeneratorConfig := policygenerator.DefaultConfig()
	newPolicyGeneratorConfig.VaultClient = newVaultClient
	newPolicyGenerator, err := policygenerator.New(newPolicyGeneratorConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	fmt.Printf("%#v\n", newPolicyGenerator)

	// TODO
}
