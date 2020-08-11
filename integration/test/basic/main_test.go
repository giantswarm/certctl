// +build k8srequired

package basic

import (
	"testing"

	"github.com/giantswarm/certctl/v2/integration/setup"
)

var (
	c setup.Config
)

func init() {
	var err error

	{
		c, err = setup.NewConfig()
		if err != nil {
			panic(err)
		}
	}
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	setup.WrapTestMain(c, m)
}
