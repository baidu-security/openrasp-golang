package goid

import (
	"bytes"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/modern-go/reflect2"
)

var expectedOffsetDict = map[string]uintptr{
	"go1.8":    192,
	"go1.8.1":  192,
	"go1.8.2":  192,
	"go1.8.3":  192,
	"go1.8.4":  192,
	"go1.8.5":  192,
	"go1.8.6":  192,
	"go1.8.7":  192,
	"go1.9":    152,
	"go1.9.1":  152,
	"go1.9.2":  152,
	"go1.10":   152,
	"go1.10.1": 152,
	"go1.10.2": 152,
	"go1.10.3": 152,
	"go1.10.4": 152,
	"go1.10.5": 152,
	"go1.10.6": 152,
	"go1.10.7": 152,
	"go1.11":   152,
	"go1.11.1": 152,
	"go1.11.2": 152,
	"go1.11.3": 152,
	"go1.11.4": 152,
}

//after go1.9 goid offset in g struct is 152, which used as default offset value
var goidOffset uintptr = 152

func init() {
	gType := reflect2.TypeByName("runtime.g").(reflect2.StructType)
	if gType == nil {
		panic("could not fetch 'runtime.g' type")
	}
	goidField := gType.FieldByName("goid")
	goidOffset = goidField.Offset()
	expectedGoidOffset, ok := expectedOffsetDict[runtime.Version()]
	if !ok {
		panic("unsupport golang version, which should >= 1.8")
	} else if goidOffset != expectedGoidOffset {
		panic("unexcepted offset of 'goid' field in 'runtime.g' struct")
	}
}

// GoIDByAsm returns the goid of current goroutine by using inline assembly
func GoIDAsm() int64 {
	g := getg()
	p_goid := (*int64)(unsafe.Pointer(g + goidOffset))
	return *p_goid
}

// GoIDByStack returns the goid of current goroutine by parsing runtime.Stack
func GoIDStack() int64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	id, _ := strconv.ParseInt(string(b), 10, 64)
	return id
}

func getg() uintptr
