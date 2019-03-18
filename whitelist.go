package openrasp

import (
	"sync"

	"github.com/baidu-security/openrasp-golang/common"
)

type WhiteList struct {
	dat common.DoubleArrayTrie
	mu  sync.RWMutex
}
