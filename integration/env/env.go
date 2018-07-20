// +build k8srequired

package env

import (
	"fmt"
	"os"
)

const (
	// EnvVaultToken is the process environment variable representing the
	// VAULT_TOKEN env var.
	EnvVaultToken = "VAULT_TOKEN"
)

var (
	vaultToken string
)

func init() {
	vaultToken = os.Getenv(EnvVaultToken)
	if vaultToken == "" {
		panic(fmt.Sprintf("env var %q must not be empty", EnvVaultToken))
	}

}

func VaultToken() string {
	return vaultToken
}
