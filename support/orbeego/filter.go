package orbeego

import (
	"net/http"

	"github.com/astaxie/beego"
	beegocontext "github.com/astaxie/beego/context"
	"github.com/baidu/openrasp"
	"github.com/baidu/openrasp/model"
	"github.com/baidu/openrasp/support/orhttp"
)

type beegoFilterStateKey struct{}

type beegoFilterState struct {
	context *beegocontext.Context
}

func init() {
	if openrasp.IsAvailable() {
		if v, err := beego.GetConfig("string", "HttpAddr", ""); err == nil {
			openrasp.GetGlobals().SetHttpAddr(v.(string))
		}
		server := model.NewServer("beego", beego.VERSION)
		openrasp.GetGlobals().SetServer(server)
	}
	AddFilters(beego.BeeApp.Handlers)
}

// Middleware returns a beego.
func Middleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return orhttp.Wrap(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(w, req)
		}))
	}
}

// AddFilters adds required filters to handlers.
func AddFilters(handlers *beego.ControllerRegister) {
	handlers.InsertFilter("*", beego.BeforeStatic, beforeStatic, false)
}

func beforeStatic(context *beegocontext.Context) {
	state, ok := context.Request.Context().Value(beegoFilterStateKey{}).(*beegoFilterState)
	if ok {
		state.context = context
	}
}
