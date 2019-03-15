package openrasp

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/baidu/openrasp/model"
	"github.com/baidu/openrasp/utils"
)

func CaptureHTTPRequestBody(req *http.Request, out *model.RequestBody) bool {
	if req.Body == nil {
		return false
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
				vcopy[i] = utils.TruncateString(v[i], GetGeneral().GetInt("body.maxbytes"))
			}
			postForm[k] = vcopy
		}
		out.Form = postForm
		return true
	}

	r := io.MultiReader(bytes.NewReader(bc.buffer.Bytes()), bc.originalBody)
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return false
	}
	out.Raw = utils.TruncateString(string(all), GetGeneral().GetInt("body.maxbytes"))
	return true
}
