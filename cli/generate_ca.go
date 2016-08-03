package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certificate-sidekick/service/vault-factory"
)

var (
	generateCACmd = &cobra.Command{
		Use:   "ca",
		Short: "Generate a certificate authority.",
		Run:   generateCARun,
	}

	generateCAVaultAddress string
	generateCAVaultToken   string
)

func init() {
	generateCmd.AddCommand(generateCACmd)

	generateCACmd.Flags().StringVar(&generateCAVaultAddress, "vault-address", "http://127.0.0.1:8200", "Address used to connect to Vault.")
	generateCACmd.Flags().StringVar(&generateCAVaultToken, "vault-token", "", "Token used to authenticate against Vault.")
}

func generateCARun(cmd *cobra.Command, args []string) {
	if generateCAVaultToken == "" {
		log.Fatalf("%#v\n", maskAnyf(invalidConfigError, "Vault token must not be empty"))
	}

	// Create Vault client and configure it with the provided admin token.
	newVaultConfig := vaultfactory.DefaultConfig()
	newVaultConfig.HTTPClient = &http.Client{}
	newVaultConfig.Address = generateCAVaultAddress
	newVaultConfig.AdminToken = generateCAVaultToken
	newVault, err := vaultfactory.New(newVaultConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	fmt.Printf("%#v\n", newVault)

	// TODO
}
