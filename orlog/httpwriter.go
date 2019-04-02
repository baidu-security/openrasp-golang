package orlog

import (
	"sync"

	"github.com/baidu-security/openrasp-golang/cloud"
)

type HttpWriter struct {
	t           string
	cm          *cloud.Client
	tokenBucket *TokenBucket
	mu          sync.Mutex
}

func NewHttpWriter(t string, cm *cloud.Client, tokenBucket *TokenBucket) *HttpWriter {
	hw := &HttpWriter{
		t:           t,
		cm:          cm,
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
	go hw.cm.Log(hw.t, p)
	return len(p), nil
}
