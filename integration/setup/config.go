// +build k8srequired

package setup

import (
	"github.com/giantswarm/e2e-harness/pkg/release"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/certctl/integration/env"
)

const (
	namespace       = "giantswarm"
	tillerNamespace = "kube-system"
)

type Config struct {
	Clients *k8sclient.Clients
	Logger  micrologger.Logger
	Release *release.Release
	Setup   *k8sclient.Setup
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

	var helmClient *helmclient.Client
	{
		c := helmclient.Config{
			Logger:    logger,
			K8sClient: cpK8sClients.K8sClient(),

			RestConfig:      cpK8sClients.RESTConfig(),
			TillerNamespace: tillerNamespace,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var newRelease *release.Release
	{
		c := release.Config{
			ExtClient:  cpK8sClients.ExtClient(),
			G8sClient:  cpK8sClients.G8sClient(),
			HelmClient: helmClient,
			K8sClient:  cpK8sClients.K8sClient(),
			Logger:     logger,

			Namespace: namespace,
		}

		newRelease, err = release.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		Clients: cpK8sClients,
		Logger:  logger,
		Release: newRelease,
		Setup:   k8sSetup,
	}

	return c, nil
}
