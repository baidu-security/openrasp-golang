package gls

import (
	"runtime"
	"sync"

	"github.com/baidu/openrasp/goid"
)

type shardItemType map[int64]map[interface{}]interface{}

var shardsCount int64 = 1 << 4

type shardPair struct {
	mu    sync.RWMutex
	items shardItemType
}

var gShards []shardPair

func (sp *shardPair) initShardItem() {
	if sp == nil {
		return
	}
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.items = make(shardItemType)
}

func (sp *shardPair) setShardItem(goid int64, value map[interface{}]interface{}) {
	if sp == nil {
		return
	}
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.items[goid] = value
}

func (sp *shardPair) getShardItem(goid int64) map[interface{}]interface{} {
	if sp == nil {
		return nil
	}
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	gls, found := sp.items[goid]
	if found {
		return gls
	} else {
		return nil
	}
}

func (sp *shardPair) removeShardItem(goid int64) {
	if sp == nil {
		return
	}
	sp.mu.Lock()
	defer sp.mu.Unlock()
	delete(sp.items, goid)
}

func init() {
	shardsCount = int64(runtime.NumCPU()) * 4
	gShards = make([]shardPair, shardsCount)
	for i := int64(0); i < shardsCount; i++ {
		gShards[i].initShardItem()
	}
}

// getGls get local storage for specified goroutine
func getGls(goid int64) map[interface{}]interface{} {
	shardIndex := goid % shardsCount
	return gShards[shardIndex].getShardItem(goid)
}

//setGls set local storage for specified goroutine
func setGls(goid int64, value map[interface{}]interface{}) {
	shardIndex := goid % shardsCount
	gShards[shardIndex].setShardItem(goid, value)
	return
}

// removeGls remove local storage for specified goroutine
func removeGls(goid int64) {
	shardIndex := goid % shardsCount
	gShards[shardIndex].removeShardItem(goid)
	return
}

// glsActivated return whether gls is activated for specified goroutine
func glsActivated(goid int64) bool {
	return getGls(goid) != nil
}

// Initialize initialize local storage for current goroutine
func Initialize() {
	setGls(goid.GoIDAsm(), make(map[interface{}]interface{}))
}

// Activated return whether gls is activated for current goroutine
func Activated() bool {
	return glsActivated(goid.GoIDAsm())
}

// Clear llear local storage for current goroutine
func Clear() {
	removeGls(goid.GoIDAsm())
}

// Get the value corresponding to a specific key from goroutine local storage
func Get(key interface{}) interface{} {
	localMap := getGls(goid.GoIDAsm())
	if localMap == nil {
		return nil
	}
	return localMap[key]
}

// Get key and value into goroutine local storage
func Set(key interface{}, value interface{}) {
	localMap := getGls(goid.GoIDAsm())
	if localMap == nil {
		panic("local storage is not activated for this goroutine")
	}
	localMap[key] = value
}
