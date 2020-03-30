package accesslog

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	"github.com/containous/traefik/v2/pkg/middlewares"
)

var (
	_ middlewares.Stateful = &CaptureResponseWriter{}
)

// captureResponseWriter is a wrapper of type http.ResponseWriter
// that tracks request status and size
type CaptureResponseWriter struct {
	Rw     http.ResponseWriter
	status int
	size   int64
}

func (crw *CaptureResponseWriter) Header() http.Header {
	return crw.Rw.Header()
}

func (crw *CaptureResponseWriter) Write(b []byte) (int, error) {
	if crw.status == 0 {
		crw.status = http.StatusOK
	}
	size, err := crw.Rw.Write(b)
	crw.size += int64(size)
	return size, err
}

func (crw *CaptureResponseWriter) WriteHeader(s int) {
	crw.Rw.WriteHeader(s)
	crw.status = s
}

func (crw *CaptureResponseWriter) Flush() {
	if f, ok := crw.Rw.(http.Flusher); ok {
		f.Flush()
	}
}

func (crw *CaptureResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := crw.Rw.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("not a hijacker: %T", crw.Rw)
}

func (crw *CaptureResponseWriter) CloseNotify() <-chan bool {
	if c, ok := crw.Rw.(http.CloseNotifier); ok {
		return c.CloseNotify()
	}
	return nil
}

func (crw *CaptureResponseWriter) Status() int {
	return crw.status
}

func (crw *CaptureResponseWriter) Size() int64 {
	return crw.size
}


