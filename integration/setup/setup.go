// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/giantswarm/appcatalog"
	"github.com/giantswarm/apptest"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/certctl/v2/integration/env"
	"github.com/giantswarm/certctl/v2/integration/key"
)

const (
	defaultTestCatalog    = "default-test"
	defaultTestCatalogURL = "https://giantswarm.github.io/default-test-catalog"
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

	name := fmt.Sprintf("%s-chart", key.VaultReleaseName())

	var latestVersion string
	{
		latestVersion, err = appcatalog.GetLatestVersion(ctx, defaultTestCatalogURL, name, "")
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var valuesYaml string
	{
		values := map[string]interface{}{
			"vault": map[string]interface{}{
				"token": env.VaultToken(),
			},
		}
		bytes, err := yaml.Marshal(values)
		if err != nil {
			return microerror.Mask(err)
		}

		valuesYaml = string(bytes)
	}

	{
		apps := []apptest.App{
			{
				CatalogName:   defaultTestCatalog,
				Name:          name,
				Namespace:     metav1.NamespaceSystem,
				Version:       latestVersion,
				ValuesYAML:    valuesYaml,
				WaitForDeploy: true,
			},
		}
		err = c.AppSetup.InstallApps(ctx, apps)
		if err != nil {
			c.Logger.LogCtx(ctx, "level", "error", "message", "install apps failed", "stack", fmt.Sprintf("%#v\n", err))
			os.Exit(-1)
		}
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
