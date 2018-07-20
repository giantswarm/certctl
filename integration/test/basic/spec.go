// +build k8srequired

package basic

const (
	defaultCertTTL        = "8640h"
	defaultCertTokenTTL   = "8640h"
	defaultClusterID      = "someid"
	defaultCommonName     = "giantswarm.io"
	defaultCertCommonName = "admin." + defaultCommonName
	defaultCATTL          = "86400h"
	defaultTokenTTL       = "720h"
)
