package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certctl/service/pki-controller"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type inspectFlags struct {
	// Vault
	VaultAddress string
	VaultToken   string

	// Cluster
	ClusterID string
}

var (
	inspectCmd = &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a Vault PKI backend including all necessary requirements.",
		Run:   inspectRun,
	}

	newInspectFlags = &inspectFlags{}
)

func init() {
	CLICmd.AddCommand(inspectCmd)

	inspectCmd.Flags().StringVar(&newInspectFlags.VaultAddress, "vault-addr", fromEnv("VAULT_ADDR", "http://127.0.0.1:8200"), "Address used to connect to Vault.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultToken, "vault-token", fromEnv("VAULT_TOKEN", ""), "Token used to authenticate against Vault.")

	inspectCmd.Flags().StringVar(&newInspectFlags.ClusterID, "cluster-id", "", "Cluster ID used to generate a new root CA for.")
}

func inspectValidate(newInspectFlags *inspectFlags) error {
	if newInspectFlags.VaultToken == "" {
		return maskAnyf(invalidConfigError, "Vault token must not be empty")
	}
	if newInspectFlags.ClusterID == "" {
		return maskAnyf(invalidConfigError, "cluster ID must not be empty")
	}

	return nil
}

func inspectRun(cmd *cobra.Command, args []string) {
	err := inspectValidate(newInspectFlags)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.HTTPClient = &http.Client{}
	newVaultFactoryConfig.Address = newInspectFlags.VaultAddress
	newVaultFactoryConfig.AdminToken = newInspectFlags.VaultToken
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

	// Create a PKI controller to setup the cluster's PKI backend including its
	// root CA and role.
	newPKIControllerConfig := pkicontroller.DefaultConfig()
	newPKIControllerConfig.VaultClient = newVaultClient
	newPKIController, err := pkicontroller.New(newPKIControllerConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	mounted, err := newPKIController.IsPKIMounted(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	generated, err := newPKIController.IsCAGenerated(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	created, err := newPKIController.IsRoleCreated(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	fmt.Printf("Inspecting cluster for ID '%s':\n", newInspectFlags.ClusterID)
	fmt.Printf("\n")
	fmt.Printf("    PKI backend mounted: %t\n", mounted)
	fmt.Printf("    Root CA generated:   %t\n", generated)
	fmt.Printf("    PKI role created:    %t\n", created)
	fmt.Printf("\n")
	fmt.Printf("Tokens might or might not be generated for this cluster.\n")
	fmt.Printf("Information about these secrets need to be inspected directly\n")
	fmt.Printf("where the cluster is installed, and cannot be retrieved through\n")
	fmt.Printf("the API.\n")
}
