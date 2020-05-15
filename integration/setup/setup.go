// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"k8s.io/helm/pkg/helm"

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
		err = c.Setup.EnsureNamespaceCreated(ctx, namespace)
		if err != nil {
			return microerror.Mask(err)
		}

		err = c.HelmClient.EnsureTillerInstalled(ctx)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var operatorTarballPath string
	{
		name := fmt.Sprintf("%s-chart", key.VaultReleaseName())
		releaseVersion, err := c.ApprClient.GetReleaseVersion(ctx, name, "stable")
		if err != nil {
			return microerror.Mask(err)
		}

		operatorTarballPath, err = c.ApprClient.PullChartTarballFromRelease(ctx, name, releaseVersion)
		if err != nil {
			return microerror.Mask(err)
		}

		c.Logger.Log("level", "debug", "message", fmt.Sprintf("tarball path '%s':", operatorTarballPath))

		defer func() {
			fs := afero.NewOsFs()
			err := fs.Remove(operatorTarballPath)
			if err != nil {
				c.Logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("deletion of %#q failed", operatorTarballPath), "stack", fmt.Sprintf("%#v", err))
			}
		}()
	}

	values := `
vault:
  token: %s
`
	values = fmt.Sprintf(values, env.VaultToken())

	err = c.HelmClient.InstallReleaseFromTarball(ctx,
		operatorTarballPath,
		namespace,
		helm.ReleaseName(key.VaultReleaseName()),
		helm.ValueOverrides([]byte(values)))
	if err != nil {
		return microerror.Mask(err)
	}

	//err = c.Release.Install(ctx, key.VaultReleaseName(), release.NewStableVersion(), values, c.Release...PodExists(ctx, "default", "app=vault"))
	//if err != nil {
	//	return microerror.Mask(err)
	//}

	{
		c.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waiting for vault pod"))

		err = c.Release.WaitForPod(ctx, namespace, "app=vault")
		if err != nil {
			return microerror.Mask(err)
		}

		c.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waited for vault pod"))
	}

	return nil
}

func teardown(c Config) error {
	return nil
}
