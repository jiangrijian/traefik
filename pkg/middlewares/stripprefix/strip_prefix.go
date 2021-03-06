package stripprefix

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	// ForwardedPrefixHeader is the default header to set prefix.
	ForwardedPrefixHeader = "X-Forwarded-Prefix"
	typeName              = "StripPrefix"
)

// stripPrefix is a middleware used to strip prefix from an URL request.
type stripPrefix struct {
	next     http.Handler
	prefixes []string
	name     string
}

// New creates a new strip prefix middleware.
func New(ctx context.Context, next http.Handler, config dynamic.StripPrefix, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")
	return &stripPrefix{
		prefixes: config.Prefixes,
		next:     next,
		name:     name,
	}, nil
}

func (s *stripPrefix) GetTracingInformation() (string, ext.SpanKindEnum) {
	return s.name, tracing.SpanKindNoneEnum
}

func (s *stripPrefix) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	fmt.Println("stripPrefix ServerHttp start()...")
	for _, prefix := range s.prefixes {
		if strings.HasPrefix(req.URL.Path, prefix) {
			fmt.Printf("before strip URL detail: %+v\n", req.URL)
			req.URL.Path = getPrefixStripped(req.URL.Path, prefix)
			if req.URL.RawPath != "" {
				req.URL.RawPath = getPrefixStripped(req.URL.RawPath, prefix)
			}
			s.serveRequest(rw, req, strings.TrimSpace(prefix))
			fmt.Printf("after strip URL detail: %+v\n", req.URL)
			return
		}
	}
	fmt.Printf("strip-middleware request detail: %+v\n", req)
	s.next.ServeHTTP(rw, req)
}

func (s *stripPrefix) serveRequest(rw http.ResponseWriter, req *http.Request, prefix string) {
	req.Header.Add(ForwardedPrefixHeader, prefix)
	req.RequestURI = req.URL.RequestURI()
	s.next.ServeHTTP(rw, req)
}

func getPrefixStripped(s, prefix string) string {
	return ensureLeadingSlash(strings.TrimPrefix(s, prefix))
}

func ensureLeadingSlash(str string) string {
	return "/" + strings.TrimPrefix(str, "/")
}
