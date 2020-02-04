package token

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var policyAlreadyExistsError = microerror.Error{
	Kind: "policyAlreadyExistsError",
}

// IsPolicyAlreadyExists asserts policyAlreadyExistsError.
func IsPolicyAlreadyExists(err error) bool {
	return microerror.Cause(err) == policyAlreadyExistsError
}
