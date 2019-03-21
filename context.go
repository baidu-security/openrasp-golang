package openrasp

import (
	"github.com/baidu-security/openrasp-golang/gls"
	"github.com/baidu-security/openrasp-golang/model"
	v8 "github.com/baidu-security/openrasp-v8/go"
)

var contextGetters *v8.ContextGetters

func DefaultContextGetters() *v8.ContextGetters {
	return contextGetters
}

func InitContextGetters() {
	contextGetters = newContextGetters()
}

func newContextGetters() *v8.ContextGetters {
	//requestInfo
	urlFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.UrlFull
	}
	pathFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.UrlPath
	}
	querystringFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.Query
	}
	methodFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.Method
	}
	protocolFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.Protocol
	}
	remoteAddrFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.RemoteAddr
	}
	parameterFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.GetBytes
	}
	jsonFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.RequestBody.Raw
	}
	bodyFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.RequestBody.Truncated
	}
	appBasePathFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.AppBasePath
	}
	headerFunc := func() interface{} {
		requestInfo, _ := gls.Get("requestInfo").(*model.RequestInfo)
		return requestInfo.HeaderBytes
	}
	//globals
	serverFunc := func() interface{} {
		return GetGlobals().ContextServer.Bytes()
	}
	cg := &v8.ContextGetters{
		Url:         urlFunc,
		Path:        pathFunc,
		Querystring: querystringFunc,
		Method:      methodFunc,
		Protocol:    protocolFunc,
		RemoteAddr:  remoteAddrFunc,
		Header:      headerFunc,
		Parameter:   parameterFunc,
		Json:        jsonFunc,
		Server:      serverFunc,
		AppBasePath: appBasePathFunc,
		Body:        bodyFunc,
	}
	return cg
}

func RequestInfoAvailable() bool {
	_, ok := gls.Get("requestInfo").(*model.RequestInfo)
	return ok
}
