// +build k8srequired

package basic

import (
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/certctl/integration/setup"
)

var (
	f *framework.Host
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var err error

	l, err := micrologger.New(micrologger.Config{})
	if err != nil {
		panic(err.Error())
	}

	f, err = framework.NewHost(framework.HostConfig{
		Logger: l,

		ClusterID: defaultClusterID,
	})
	if err != nil {
		panic(err.Error())
	}

	setup.WrapTestMain(f, m)
}
