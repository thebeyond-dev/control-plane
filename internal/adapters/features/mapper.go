package features

import (
	unleashcontext "github.com/Unleash/unleash-go-sdk/v6/context"
	"github.com/thebeyond-net/control-plane/internal/core/application/input"
)

func toUnleashContext(ctx input.FeatureContext) unleashcontext.Context {
	return unleashcontext.Context{
		UserId:        ctx.UserID,
		SessionId:     ctx.SessionID,
		RemoteAddress: ctx.RemoteAddress,
		Environment:   ctx.Environment,
		AppName:       ctx.AppName,
		CurrentTime:   ctx.CurrentTime,
		Properties:    ctx.Properties,
	}
}
