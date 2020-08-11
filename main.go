package main

import (
	"log"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/certctl/v2/cli"
)

func main() {
	if err := cli.CLICmd.Execute(); err != nil {
		log.Fatalf("%#v\n", microerror.Mask(err))
	}
}
