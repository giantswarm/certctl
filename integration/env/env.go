// +build k8srequired

package env

import (
	"fmt"
	"os"
)

const (
	EnvVarE2EKubeconfig = "E2E_KUBECONFIG"
	// EnvVaultToken is the process environment variable representing the
	// VAULT_TOKEN env var.
	EnvVaultToken = "VAULT_TOKEN"
)

var (
	kubeconfig string
	vaultToken string
)

func init() {
	kubeconfig = os.Getenv(EnvVarE2EKubeconfig)
	if kubeconfig == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarE2EKubeconfig))
	}

	vaultToken = os.Getenv(EnvVaultToken)
	if vaultToken == "" {
		vaultToken = "myToken"
	}
}

func KubeConfigPath() string {
	return kubeconfig
}

func VaultToken() string {
	return vaultToken
}
