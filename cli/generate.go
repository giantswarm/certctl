package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate certificates.",
		Run:   generateRun,
	}
)

func init() {
	CLICmd.AddCommand(generateCmd)
}

func generateRun(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
	os.Exit(1)
}
