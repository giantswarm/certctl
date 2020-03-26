// +build k8srequired

package setup

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/release"
	"github.com/giantswarm/e2etemplates/pkg/chartvalues"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/certctl/integration/env"
	"github.com/giantswarm/certctl/integration/key"
)

func WrapTestMain(c Config, m *testing.M) {
	var v int
	var err error

	err = setup(c)
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	if v == 0 {
		v = m.Run()
	}

	err = teardown(c)
	if err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	os.Exit(v)
}

func setup(c Config) error {
	var err error

	ctx := context.Background()

	{
		err = c.HelmClient.EnsureTillerInstalled(ctx)
		if err != nil {
			return microerror.Mask(err)
		}

		err = c.Setup.EnsureNamespaceCreated(ctx, namespace)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var values string
	var err error
	{
		c := chartvalues.E2ESetupVaultConfig{
			Vault: chartvalues.E2ESetupVaultConfigVault{
				Token: env.VaultToken(),
			},
		}

		values, err = chartvalues.NewE2ESetupVault(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	err = c.Release.Install(ctx, key.VaultReleaseName(), release.NewStableVersion(), values, c.Release.Condition().PodExists(ctx, "default", "app=vault"))
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func teardown(c Config) error {
	return nil
}
