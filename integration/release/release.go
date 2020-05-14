// +build k8srequired

package release

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	HelmClient helmclient.Interface
	K8sClient  k8sclient.Interface
	Logger     micrologger.Logger
}

type Release struct {
	helmClient helmclient.Interface
	k8sClient  k8sclient.Interface
	logger     micrologger.Logger
}

func New(config Config) (*Release, error) {
	if config.HelmClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HelmClient must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Release{
		helmClient: config.HelmClient,
		k8sClient:  config.K8sClient,
		logger:     config.Logger,
	}

	return r, nil
}

func (r *Release) WaitForPod(ctx context.Context, namespace, labelSelector string) error {
	o := func() error {
		pods, err := r.k8sClient.K8sClient().CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return microerror.Mask(err)
		}
		if len(pods.Items) != 1 {
			return microerror.Maskf(waitError, "expected 1 pod but got %d", len(pods.Items))
		}

		pod := pods.Items[0]
		if pod.Status.Phase != corev1.PodRunning {
			return microerror.Maskf(waitError, "expected Pod phase %#q but got %#q", corev1.PodRunning, pod.Status.Phase)
		}

		return nil
	}

	n := func(err error, t time.Duration) {
		r.logger.Log("level", "debug", "message", fmt.Sprintf("failed to get pod with selector '%s': retrying in %s", labelSelector, t), "stack", fmt.Sprintf("%v", err))
	}

	b := backoff.NewExponential(10*time.Minute, 60*time.Second)
	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
