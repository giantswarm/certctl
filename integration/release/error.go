//go:build k8srequired

package release

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var releaseStatusNotMatchingError = &microerror.Error{
	Kind: "releaseStatusNotMatchingError",
}

var releaseVersionNotMatchingError = &microerror.Error{
	Kind: "releaseVersionNotMatchingError",
}

// IsReleaseVersionNotMatching asserts releaseVersionNotMatchingError
func IsReleaseVersionNotMatching(err error) bool {
	return microerror.Cause(err) == releaseVersionNotMatchingError
}

var waitError = &microerror.Error{
	Kind: "waitError",
}

// IsWait asserts waitError.
func IsWait(err error) bool {
	return microerror.Cause(err) == waitError
}
