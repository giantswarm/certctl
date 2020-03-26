// +build k8srequired

package env

import (
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
		vaultToken = "myToken"
	}
}

func VaultToken() string {
	return vaultToken
}
