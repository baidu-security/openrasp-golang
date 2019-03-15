package model

import (
	"net/http"
	"net/url"
	"runtime"

	"github.com/baidu/openrasp/utils"
)

type ContextInfo struct {
	AppBasePath string         `json:"appBasePath"`
	RemoteAddr  string         `json:"remoteAddr"`
	Protocol    string         `json:"protocol"`
	Method      string         `json:"method"`
	Query       string         `json:"querystring"`
	Path        string         `json:"path"`
	UrlFull     string         `json:"url"`
	Get         url.Values     `json:"get"`
	Header      http.Header    `json:"header"`
	Server      *ContextServer `json:"server"`
	*RequestBody
}

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

func NewContextInfo(request *http.Request, cs *ContextServer, rb *RequestBody) *ContextInfo {
	ci := &ContextInfo{
		AppBasePath: "",
		RemoteAddr:  request.RemoteAddr,
		Protocol:    request.Proto,
		Method:      request.Method,
		Query:       request.URL.RawQuery,
		UrlFull:     request.URL.String(),
		Path:        request.URL.Path,
		Get:         request.URL.Query(),
		Header:      request.Header,
		Server:      cs,
	}
	ci.RequestBody = rb
	return ci
}
