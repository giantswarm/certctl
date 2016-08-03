package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	CLICmd = &cobra.Command{
		Use:   "certificate-sidekick",
		Short: "A sidekick process able to request certificate generation from Vault to write files to the local filesystem.",

		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, nil)
			os.Exit(1)
		},
	}
)
