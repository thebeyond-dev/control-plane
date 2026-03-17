package features

import (
	"github.com/Unleash/unleash-go-sdk/v6"
	"github.com/thebeyond-net/control-plane/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type unleashAdapter struct{}

func NewUnleashAdapter() ports.FeatureFlags {
	return &unleashAdapter{}
}

func (a *unleashAdapter) IsEnabled(name string, ctx input.FeatureContext) bool {
	return unleash.IsEnabled(name, unleash.FeatureOptions{
		Ctx: toUnleashContext(ctx),
	})
}
