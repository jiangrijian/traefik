package chain

import (
	"context"
	"net/http"

	"github.com/containous/alice"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
)

const (
	typeName = "Chain"
)

type chainBuilder interface {
	BuildChain(ctx context.Context, middlewares []string, service string) *alice.Chain
}

// New creates a chain middleware
func New(ctx context.Context, next http.Handler, config dynamic.Chain, builder chainBuilder, name string, service string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")

	middlewareChain := builder.BuildChain(ctx, config.Middlewares, service)
	return middlewareChain.Then(next)
}
