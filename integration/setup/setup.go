// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/giantswarm/helmclient/v2"
	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/certctl/v2/integration/env"
	"github.com/giantswarm/certctl/v2/integration/key"
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

	values := map[string]interface{}{
		"vault": map[string]interface{}{
			"token": env.VaultToken(),
		},
	}

	opts := helmclient.InstallOptions{
		ReleaseName: key.VaultReleaseName(),
	}

	err = c.HelmClient.InstallReleaseFromTarball(ctx,
		operatorTarballPath,
		namespace,
		values,
		opts)
	if err != nil {
		return microerror.Mask(err)
	}

	{
		c.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("waiting for vault pod"))

		err = c.Release.WaitForPod(ctx, metav1.NamespaceDefault, "app=vault")
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
