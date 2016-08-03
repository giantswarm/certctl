package main

import (
	"log"
	"os"

	"github.com/giantswarm/certificate-sidekick/cli"
)

func main() {
	if err := cli.CLICmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
