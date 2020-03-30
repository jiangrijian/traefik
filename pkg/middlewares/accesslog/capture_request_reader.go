package accesslog

import "io"

type CaptureRequestReader struct {
	Source io.ReadCloser
	Count  int64
}

func (r *CaptureRequestReader) Read(p []byte) (int, error) {
	n, err := r.Source.Read(p)
	r.Count += int64(n)
	return n, err
}

func (r *CaptureRequestReader) Close() error {
	return r.Source.Close()
}
