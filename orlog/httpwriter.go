package orlog

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type HttpWriter struct {
	url         string
	appId       string
	appSecret   string
	client      *http.Client
	tokenBucket *TokenBucket
	mu          sync.Mutex
}

func NewHttpWriter(url, appId, appSecret string, tokenBucket *TokenBucket) *HttpWriter {
	tr := &http.Transport{
		IdleConnTimeout:    20 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	hw := &HttpWriter{
		url:         url,
		appId:       appId,
		appSecret:   appSecret,
		client:      client,
		tokenBucket: tokenBucket,
	}
	return hw
}

func (hw *HttpWriter) Write(p []byte) (n int, err error) {
	hw.mu.Lock()
	defer hw.mu.Unlock()
	if hw.tokenBucket != nil && hw.tokenBucket.Consume() {
		return 0, nil
	}
	req, err := http.NewRequest("POST", hw.url, bytes.NewReader(p))
	req.Header.Add("X-OpenRASP-AppID", hw.appId)
	req.Header.Add("X-OpenRASP-AppSecret", hw.appSecret)
	resp, err := hw.client.Do(req)
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	return len(p), err
}
