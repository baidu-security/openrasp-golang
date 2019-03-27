package orhttp

import (
	"net/http"

	openrasp "github.com/baidu-security/openrasp-golang"
	"github.com/baidu-security/openrasp-golang/gls"
	"github.com/baidu-security/openrasp-golang/model"
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
	if openrasp.IsComplete() {
		gls.Initialize()
		defer func() {
			gls.Clear()
		}()
		whiteUrl := openrasp.ExtractWhiteKey(req.URL)
		whiteBitMask := openrasp.GetWhite().PrefixSearch(whiteUrl)
		gls.Set("whiteMask", whiteBitMask)

		clientIpHeader := openrasp.GetGeneral().GetString("clientip.header")
		bodyMaxByte := openrasp.GetGeneral().GetInt("body.maxbytes")
		requestInfo := model.NewRequestInfo(req, clientIpHeader, bodyMaxByte)
		gls.Set("requestInfo", requestInfo)

		w, resp := WrapResponseWriter(w, req)
		gls.Set("responseWriter", w)
		defer func() {
			if v := recover(); v != nil {
				if resp.StatusCode == 0 {
					w.WriteHeader(http.StatusOK)
				}
			}
		}()
	}
	h.handler.ServeHTTP(w, req)
}
