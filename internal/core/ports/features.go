package ports

import "github.com/thebeyond-net/control-plane/internal/core/application/input"

type FeatureFlags interface {
	IsEnabled(name string, ctx input.FeatureContext) bool
}
