package model

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/baidu/openrasp/utils"
)

type RequestInfo struct {
	Method       string      `json:"request_method"`
	UrlFull      string      `json:"url"`
	UrlHost      string      `json:"target"`
	UrlPath      string      `json:"path"`
	AttackSource string      `json:"attack_source"`
	ClientIp     string      `json:"client_ip"`
	RequestId    string      `json:"request_id"`
	Header       http.Header `json:"header"`
	*RequestBody
}

type RequestBody struct {
	Raw  string     `json:"body"`
	Form url.Values `json:"form"`
}

func NewRequestInfo(request *http.Request, clientIpHeader string, bodySize int) *RequestInfo {
	rb := NewRequestBody(request, bodySize)
	ri := &RequestInfo{
		Method:       request.Method,
		UrlFull:      request.URL.String(),
		UrlHost:      request.URL.Host,
		UrlPath:      request.URL.Path,
		AttackSource: request.RemoteAddr,
		ClientIp:     request.Header.Get(clientIpHeader),
		RequestId:    utils.GenerateRequestId(),
	}
	ri.SetRequestBody(rb)
	return ri
}

func (ri *RequestInfo) SetRequestBody(rb *RequestBody) {
	ri.RequestBody = rb
}

func NewRequestBody(req *http.Request, size int) *RequestBody {
	out := &RequestBody{}

	if req.Body == nil {
		return out
	}

	type bodyCapturer struct {
		originalBody io.ReadCloser
		buffer       bytes.Buffer
		request      *http.Request
	}

	type readerCloser struct {
		io.Reader
		io.Closer
	}

	bc := bodyCapturer{
		request:      req,
		originalBody: req.Body,
	}

	req.Body = &readerCloser{
		Reader: io.TeeReader(req.Body, &bc.buffer),
		Closer: req.Body,
	}
	if bc.request.PostForm != nil {
		postForm := make(url.Values, len(bc.request.PostForm))
		for k, v := range bc.request.PostForm {
			vcopy := make([]string, len(v))
			for i := range vcopy {
				vcopy[i] = utils.TruncateString(v[i], size)
			}
			postForm[k] = vcopy
		}
		out.Form = postForm
		return out
	}

	r := io.MultiReader(bytes.NewReader(bc.buffer.Bytes()), bc.originalBody)
	all, err := ioutil.ReadAll(r)
	if err == nil {
		out.Raw = utils.TruncateString(string(all), size)
	}
	return out
}
