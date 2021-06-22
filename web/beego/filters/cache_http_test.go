package filters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/linmadan/egglib-go/persistent/cache/gzcache"
	"github.com/stretchr/testify/assert"
)

func TestAtomicPersistenceQueryHandler(t *testing.T) {
	nodeCache := gzcache.NewNodeCache("127.0.0.1:6379", "")
	apq := AtomicPersistenceQueryHandler(
		WithCache(nodeCache),
		WithServiceName("testapq"),
		WithExpire(60*60), //after 60min object will expire
		WithRouters([]string{
			"/ok",
		}),
	)

	server := web.NewHttpSever()
	server.Get("/ok", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilterChain("/*", apq)
	r := httptest.NewRequest(http.MethodGet, "/ok", nil)
	SecureHttpRequest(r, "")
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, "ok", w.Body.String())
}
