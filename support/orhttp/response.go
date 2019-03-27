package orhttp

import (
	"net/http"
	"strings"

	openrasp "github.com/baidu-security/openrasp-golang"
	"github.com/baidu-security/openrasp-golang/gls"
	"github.com/baidu-security/openrasp-golang/model"
)

type OpenRASPBlocker interface {
	BlockByOpenRASP()
}

func WrapResponseWriter(w http.ResponseWriter, req *http.Request) (http.ResponseWriter, *Response) {
	rw := ResponseWriter{
		ResponseWriter: w,
		resp: Response{
			Headers: w.Header(),
			req:     req,
		},
	}
	h, _ := w.(http.Hijacker)
	p, _ := w.(http.Pusher)
	switch {
	case h != nil && p != nil:
		rwhp := &responseWriterHijackerPusher{
			ResponseWriter: rw,
			Hijacker:       h,
			Pusher:         p,
		}
		return rwhp, &rwhp.resp
	case h != nil:
		rwh := &responseWriterHijacker{
			ResponseWriter: rw,
			Hijacker:       h,
		}
		return rwh, &rwh.resp
	case p != nil:
		rwp := &responseWriterPusher{
			ResponseWriter: rw,
			Pusher:         p,
		}
		return rwp, &rwp.resp
	}
	return &rw, &rw.resp
}

type Response struct {
	StatusCode int
	Sent       bool
	Headers    http.Header
	req        *http.Request
}

type ResponseWriter struct {
	http.ResponseWriter
	resp Response
}

func (res *Response) detectContentType() string {
	ct := res.Headers.Get("Content-Type")
	if len(ct) > 0 {
		return ct
	} else {
		return res.req.Header.Get("Accept")
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	requestInfo, ok := gls.Get("requestInfo").(*model.RequestInfo)
	if ok {
		w.ResponseWriter.Header().Set("X-Request-ID", requestInfo.GetRequestId())
	}
	w.ResponseWriter.Header().Set("X-Protected-By", "OpenRASP")
	w.ResponseWriter.WriteHeader(statusCode)
	w.resp.StatusCode = statusCode
	w.resp.Sent = true
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if w.resp.StatusCode == 0 {
		w.resp.StatusCode = http.StatusOK
	}
	return n, err
}

func (w *ResponseWriter) CloseNotify() <-chan bool {
	if closeNotifier, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return closeNotifier.CloseNotify()
	}
	return nil
}

func (w *ResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *ResponseWriter) HasSent() bool {
	return w.resp.Sent
}

func (w *ResponseWriter) appendBlockContent(requestId string) {
	ct := w.resp.detectContentType()
	if strings.HasPrefix(ct, "application/json") {
		contentJson := openrasp.GetGeneral().GetString("block.content_json")
		contentJson = strings.Replace(contentJson, "%request_id%", requestId, -1)
		w.ResponseWriter.Write([]byte(contentJson))
	} else if strings.HasPrefix(ct, "application/xml") || strings.HasPrefix(ct, "text/xml") {
		contentXml := openrasp.GetGeneral().GetString("block.content_xml")
		contentXml = strings.Replace(contentXml, "%request_id%", requestId, -1)
		w.ResponseWriter.Write([]byte(contentXml))
	} else {
		contentHtml := openrasp.GetGeneral().GetString("block.content_html")
		contentHtml = strings.Replace(contentHtml, "%request_id%", requestId, -1)
		w.ResponseWriter.Write([]byte(contentHtml))
	}
}

func (w *ResponseWriter) BlockByOpenRASP() {
	var requestId string
	requestInfo, ok := gls.Get("requestInfo").(*model.RequestInfo)
	if ok {
		requestId = requestInfo.GetRequestId()
	}
	if !w.resp.Sent {
		redirectUrl := openrasp.GetGeneral().GetString("block.redirect_url")
		redirectUrl = strings.Replace(redirectUrl, "%request_id%", requestId, -1)
		w.ResponseWriter.Header().Set("Location", redirectUrl)
		statusCode := openrasp.GetGeneral().GetInt("block.status_code")
		w.WriteHeader(statusCode)
	} else {
		w.appendBlockContent(requestId)
	}
	panic(openrasp.ErrBlock)
}

type responseWriterHijacker struct {
	ResponseWriter
	http.Hijacker
}

type responseWriterPusher struct {
	ResponseWriter
	http.Pusher
}

type responseWriterHijackerPusher struct {
	ResponseWriter
	http.Hijacker
	http.Pusher
}
