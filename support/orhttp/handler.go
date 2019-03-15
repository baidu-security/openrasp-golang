package orhttp

import (
	"net/http"

	"github.com/baidu/openrasp/gls"
)

// Wrap returns an http.Handler wrapping h
func Wrap(h http.Handler) http.Handler {
	if h == nil {
		panic("h == nil")
	}
	handler := &handler{
		handler: h,
	}
	return handler
}

// handler wraps an http.Handler
type handler struct {
	handler http.Handler
}

// ServeHTTP delegates to h.Handler
func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	gls.Initialize()
	defer func() {
		gls.Clear()
	}()
	gls.Set("request", req)
	w, resp := WrapResponseWriter(w)
	defer func() {
		if v := recover(); v != nil {
			if resp.StatusCode == 0 {
				w.WriteHeader(http.StatusOK)
			}
		}
	}()
	h.handler.ServeHTTP(w, req)
}

// WrapResponseWriter wraps an http.ResponseWriter and returns the wrapped
// value along with a *Response which will be filled in when the handler
// is called. The *Response value must not be inspected until after the
// request has been handled, to avoid data races. If neither of the
// ResponseWriter's Write or WriteHeader methods are called, then the
// response's StatusCode field will be zero.
//
// The returned http.ResponseWriter implements http.Pusher and http.Hijacker
// if and only if the provided http.ResponseWriter does.
func WrapResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *Response) {
	rw := responseWriter{
		ResponseWriter: w,
		resp: Response{
			Headers: w.Header(),
		},
	}
	h, _ := w.(http.Hijacker)
	p, _ := w.(http.Pusher)
	switch {
	case h != nil && p != nil:
		rwhp := &responseWriterHijackerPusher{
			responseWriter: rw,
			Hijacker:       h,
			Pusher:         p,
		}
		return rwhp, &rwhp.resp
	case h != nil:
		rwh := &responseWriterHijacker{
			responseWriter: rw,
			Hijacker:       h,
		}
		return rwh, &rwh.resp
	case p != nil:
		rwp := &responseWriterPusher{
			responseWriter: rw,
			Pusher:         p,
		}
		return rwp, &rwp.resp
	}
	return &rw, &rw.resp
}

// Response records details of the HTTP response.
type Response struct {
	// StatusCode records the HTTP status code set via WriteHeader.
	StatusCode int

	// Headers holds the headers set in the ResponseWriter.
	Headers http.Header
}

type responseWriter struct {
	http.ResponseWriter
	resp Response
}

func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader sets w.resp.StatusCode and calls through to the embedded
// ResponseWriter.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.resp.StatusCode = statusCode
}

// Write calls through to the embedded ResponseWriter, setting
// w.resp.StatusCode to http.StatusOK if WriteHeader has not already
// been called.
func (w *responseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if w.resp.StatusCode == 0 {
		w.resp.StatusCode = http.StatusOK
	}
	return n, err
}

// CloseNotify returns w.closeNotify() if w.closeNotify is non-nil,
// otherwise it returns nil.
func (w *responseWriter) CloseNotify() <-chan bool {
	if closeNotifier, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return closeNotifier.CloseNotify()
	}
	return nil
}

// Flush calls w.flush() if w.flush is non-nil, otherwise
// it does nothing.
func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

type responseWriterHijacker struct {
	responseWriter
	http.Hijacker
}

type responseWriterPusher struct {
	responseWriter
	http.Pusher
}

type responseWriterHijackerPusher struct {
	responseWriter
	http.Hijacker
	http.Pusher
}
