package openrasp

import (
	"net/url"
	"sort"
	"sync"

	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/gls"
)

type WhiteList struct {
	dat *common.DoubleArrayTrie
	mu  sync.RWMutex
}

func NewWhiteList() *WhiteList {
	wl := &WhiteList{
		dat: common.NewDoubleArrayTrie(),
	}
	allWhiteMap := map[string]int{
		"": common.AllType,
	}
	wl.Build(allWhiteMap)
	return wl
}

func (wl *WhiteList) build(urls []string, bits []int) {
	wl.mu.Lock()
	defer wl.mu.Unlock()
	wl.dat.Clear()
	wl.dat.Build(urls, nil, bits, len(urls))
}

func (wl *WhiteList) Build(m map[string]int) {
	var urls []string
	var bits []int
	for k := range m {
		urls = append(urls, k)
	}
	sort.Strings(urls)
	for _, url := range urls {
		bits = append(bits, m[url])
	}
	wl.build(urls, bits)
}

func (wl *WhiteList) PrefixSearch(key string) int {
	result := 0
	bitMasks := wl.prefixSearch(key)
	for _, v := range bitMasks {
		result |= v
	}
	return result
}

func (wl *WhiteList) prefixSearch(key string) []int {
	wl.mu.RLock()
	defer wl.mu.RUnlock()
	return wl.dat.CommonPrefixSearch(key)
}

func (wl *WhiteList) UpdateWhiteList() {
	hookWhite := GetGeneral().GetStringMap("hook.white")
	codeWhiteMap := make(map[string]int, len(hookWhite))
	for k, v := range hookWhite {
		mask := 0
		typeLists, ok := v.([]interface{})
		if ok {
			for _, typeValue := range typeLists {
				typeString, ok := typeValue.(string)
				if ok {
					checkType := common.CheckStringToType(typeString)
					mask |= int(checkType)
				}
			}
			if mask != 0 {
				if k == "*" {
					k = ""
				}
				codeWhiteMap[k] = mask
			}
		}
	}
	wl.Build(codeWhiteMap)
}

func (wl *WhiteList) OnConfigUpdate() {
	wl.UpdateWhiteList()
}

func ExtractWhiteKey(u *url.URL) string {
	fullUrl := u.String()
	if u.Scheme != "" {
		fullUrl = fullUrl[(len(u.Scheme) + 1):]
	}
	if u.Opaque != "" {
		return fullUrl
	} else {
		if u.Scheme != "" || u.Host != "" || u.User != nil {
			if u.Host != "" || u.Path != "" || u.User != nil {
				fullUrl = fullUrl[len("//"):]
			}
		}
		return fullUrl
	}
}

func WhitelistOption(ac common.AttackChecker) bool {
	bitMaskValue := gls.Get("whiteMask")
	bitMask, ok := bitMaskValue.(int)
	if ok && (bitMask&int(ac.GetType()) != 0) {
		return true
	}
	return false
}
