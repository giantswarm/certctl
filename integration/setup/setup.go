// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/giantswarm/appcatalog"
	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"k8s.io/helm/pkg/helm"

	"github.com/giantswarm/certctl/integration/key"
)

const (
	vaultAppRelease = "vault-app"
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

	var tarballPath string
	{
		tarballURL, err := appcatalog.GetLatestChart(ctx, key.OperationPlatformCatalogStorageURL(), vaultAppRelease, "")
		if err != nil {
			return microerror.Mask(err)
		}

		tarballPath, err = c.HelmClient.PullChartTarball(ctx, tarballURL)
		if err != nil {
			return microerror.Mask(err)
		}

		defer func() {
			fs := afero.NewOsFs()
			err := fs.Remove(tarballPath)
			if err != nil {
				c.Logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("deletion of %#q failed", tarballPath), "stack", fmt.Sprintf("%#v", err))
			}
		}()
	}

	// VaultAppValues values required by chart-operator-chart.
	const VaultAppValues = `
namespace: "giantswarm"
storage:
  size: 512Mi
`

	err = c.HelmClient.InstallReleaseFromTarball(ctx,
		tarballPath,
		namespace,
		helm.ReleaseName(key.VaultReleaseName()),
		helm.ValueOverrides([]byte(VaultAppValues)))
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
