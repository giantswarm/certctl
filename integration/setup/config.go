//go:build k8srequired

package setup

import (
	"github.com/giantswarm/apptest"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/certctl/v2/integration/env"
	"github.com/giantswarm/certctl/v2/integration/release"
)

const (
	namespace = "giantswarm"
)

type Config struct {
	AppSetup *apptest.AppSetup
	Clients  k8sclient.Interface
	Logger   micrologger.Logger
	Release  *release.Release
	Setup    *k8sclient.Setup
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

	var appSetup *apptest.AppSetup
	{
		c := apptest.Config{
			Logger: logger,

			KubeConfigPath: env.KubeConfigPath(),
		}
		appSetup, err = apptest.New(c)
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
		AppSetup: appSetup,
		Clients:  cpK8sClients,
		Logger:   logger,
		Release:  newRelease,
		Setup:    k8sSetup,
	}

	return c, nil
}
