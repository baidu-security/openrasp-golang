package openrasp

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/baidu-security/openrasp-golang/common"
	"github.com/baidu-security/openrasp-golang/model"
	v8 "github.com/baidu-security/openrasp-v8/go"
)

type BuildinAction struct {
	actionMap map[common.CheckType]model.InterceptCode
	mu        sync.RWMutex
}

func NewBuildinAction() *BuildinAction {
	return &BuildinAction{
		actionMap: make(map[common.CheckType]model.InterceptCode),
	}
}

func (ba *BuildinAction) Clear() {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.actionMap = nil
}

func (ba *BuildinAction) Set(ct common.CheckType, ic model.InterceptCode) {
	ba.mu.Lock()
	defer ba.mu.Unlock()
	ba.actionMap[ct] = ic
}

func (ba *BuildinAction) Get(ct common.CheckType) model.InterceptCode {
	ba.mu.RLock()
	defer ba.mu.RUnlock()
	ic, ok := ba.actionMap[ct]
	if ok {
		return ic
	}
	return model.Ignore
}

func (ba *BuildinAction) OnPluginUpdate() {
	script := common.BuildinActionScript()
	if len(script) > 0 {
		actionResult := v8.ExecScript(script, "extract_buildin_action")
		var ms [][]string
		unquotedResult, err := strconv.Unquote(actionResult)
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(unquotedResult), &ms)
		if err == nil {
			for _, stringPair := range ms {
				if len(stringPair) == 2 {
					ba.Set(common.CheckStringToType(stringPair[0]), model.InterceptStringToCode(stringPair[1]))
				}
			}
		}
	}
}
