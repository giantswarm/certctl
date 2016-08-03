package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certificate-sidekick/service/vault-factory"
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

	// Create Vault client and configure it with the provided admin token.
	newVaultConfig := vaultfactory.DefaultConfig()
	newVaultConfig.HTTPClient = &http.Client{}
	newVaultConfig.Address = generateSignedVaultAddress
	newVaultConfig.AdminToken = generateSignedVaultToken
	newVault, err := vaultfactory.New(newVaultConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	fmt.Printf("%#v\n", newVault)

	// TODO
}
