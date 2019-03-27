package model

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/baidu-security/openrasp-golang/utils"
)

type ContextServer struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	OS       string `json:"os"`
}

func NewContextServer() *ContextServer {
	cs := &ContextServer{
		Language: "golang",
		Name:     "GO",
		Version:  runtime.Version(),
		OS:       utils.GetOs(),
	}
	return cs
}

func (cs *ContextServer) Bytes() []byte {
	b, err := json.Marshal(cs)
	if err != nil {
		return nil
	} else {
		return b
	}
}

type RequestInfo struct {
	Method       string            `json:"request_method"`
	UrlFull      string            `json:"url"`
	UrlHost      string            `json:"target"`
	UrlPath      string            `json:"path"`
	AttackSource string            `json:"attack_source"`
	ClientIp     string            `json:"client_ip"`
	RequestId    string            `json:"request_id"`
	Header       map[string]string `json:"header"`
	Query        string            `json:"-"`
	Protocol     string            `json:"-"`
	RemoteAddr   string            `json:"-"`
	Get          map[string]string `json:"-"`
	AppBasePath  string            `json:"-"`
	HeaderBytes  []byte            `json:"-"`
	GetBytes     []byte            `json:"-"`
	*RequestBody
}

type RequestBody struct {
	Raw       string     `json:"-"`
	Truncated string     `json:"body"`
	Form      url.Values `json:"form"`
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
		Header:       joinHeader(request.Header),
		RequestId:    utils.GenerateRequestId(),
		Query:        request.URL.RawQuery,
		Protocol:     request.Proto,
		RemoteAddr:   request.RemoteAddr,
		Get:          extractHeader(request.URL.Query()),
		AppBasePath:  "",
	}
	headerBytes, _ := json.Marshal(ri.Header)
	ri.HeaderBytes = headerBytes

	getBytes, _ := json.Marshal(ri.Get)
	ri.GetBytes = getBytes

	ri.SetRequestBody(rb)
	return ri
}

func joinHeader(source map[string][]string) map[string]string {
	joinMap := make(map[string]string, len(source))
	for k, headers := range source {
		joinMap[k] = strings.Join(headers, ", ")
	}
	return joinMap
}

func extractHeader(source map[string][]string) map[string]string {
	extractMap := make(map[string]string, len(source))
	for k, _ := range source {
		extractMap[k] = http.Header(source).Get(k)
	}
	return extractMap
}

func (ri *RequestInfo) SetRequestBody(rb *RequestBody) {
	ri.RequestBody = rb
}

func (ri *RequestInfo) GetRequestId() string {
	return ri.RequestId
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
				vcopy[i] = v[i]
			}
			postForm[k] = vcopy
		}
		out.Form = postForm
		return out
	}

	r := io.MultiReader(bytes.NewReader(bc.buffer.Bytes()), bc.originalBody)
	all, err := ioutil.ReadAll(r)
	if err == nil {
		out.Raw = string(all)
		out.Truncated = utils.TruncateString(string(all), size)
	}
	return out
}
