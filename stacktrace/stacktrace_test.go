package stacktrace

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStacktrace(t *testing.T) {
	expect := []string{
		"github.com/baidu/openrasp/stacktrace.TestStacktrace.func1",
		"runtime.call32",
		"runtime.gopanic",
		"github.com/baidu/openrasp/stacktrace.(*panicker).panic",
		"github.com/baidu/openrasp/stacktrace.TestStacktrace",
	}
	defer func() {
		err := recover()
		if err == nil {
			t.FailNow()
		}
		allFrames := AppendStacktrace(nil, 1, 5)
		functions := make([]string, len(allFrames))
		for i, frame := range allFrames {
			functions[i] = frame.Function
		}
		if diff := cmp.Diff(functions, expect); diff != "" {
			t.Fatalf("%s", diff)
		}
	}()
	(&panicker{}).panic()
}

func TestLogStacktrace(t *testing.T) {
	expect := []string{
		"/github.com/baidu/openrasp/stacktrace/stacktrace_test.go(github.com/baidu/openrasp/stacktrace.TestStacktrace.func1:47)",
		"/runtime/asm_amd64.s(runtime.call32:522)",
		"/src/runtime/panic.go(runtime.gopanic:513)",
		"/github.com/baidu/openrasp/stacktrace/stacktrace_test.go(github.com/baidu/openrasp/stacktrace.(*panicker).panic:67)",
		"/github.com/baidu/openrasp/stacktrace/stacktrace_test.go(github.com/baidu/openrasp/stacktrace.TestLogStacktrace:61)",
	}
	defer func() {
		err := recover()
		if err == nil {
			t.FailNow()
		}
		allFrames := AppendStacktrace(nil, 1, 5)
		logStacks := LogFormat(allFrames)
		expectLen := len(expect)
		if len(logStacks) != expectLen {
			t.Fatalf("different stack depth")
		} else {
			for i := 1; i < expectLen; i++ {
				if !strings.HasSuffix(logStacks[i], expect[i]) {
					t.Fatalf("%s dose not end with %s", expect[i], logStacks[i])
				}
			}
		}
	}()
	(&panicker{}).panic()
}

type panicker struct{}

func (*panicker) panic() {
	panic("oh noes")
}
