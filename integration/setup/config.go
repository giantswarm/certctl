// +build k8srequired

package setup

import (
	"github.com/giantswarm/apprclient"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/certctl/integration/env"
	"github.com/giantswarm/certctl/integration/release"
)

const (
	namespace = "giantswarm"
)

type Config struct {
	ApprClient apprclient.Interface
	HelmClient helmclient.Interface
	Clients    k8sclient.Interface
	Logger     micrologger.Logger
	Release    *release.Release
	Setup      *k8sclient.Setup
}

func NewConfig() (Config, error) {
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var apprClient apprclient.Interface
	{
		c := apprclient.Config{
			Logger: logger,

			Address:      "https://quay.io",
			Organization: "giantswarm",
		}

		apprClient, err = apprclient.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var cpK8sClients *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger: logger,

			KubeConfigPath: env.KubeConfigPath(),
		}

		cpK8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var k8sSetup *k8sclient.Setup
	{
		c := k8sclient.SetupConfig{
			Clients: cpK8sClients,
			Logger:  logger,
		}

		k8sSetup, err = k8sclient.NewSetup(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var helmClient helmclient.Interface
	{
		c := helmclient.Config{
			Logger:    logger,
			K8sClient: cpK8sClients,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var newRelease *release.Release
	{
		c := release.Config{
			K8sClient: cpK8sClients,
			Logger:    logger,
		}

		newRelease, err = release.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		ApprClient: apprClient,
		Clients:    cpK8sClients,
		Logger:     logger,
		HelmClient: helmClient,
		Release:    newRelease,
		Setup:      k8sSetup,
	}

	return c, nil
}
