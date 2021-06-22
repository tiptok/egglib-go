package filters

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"

	// "github.com/tiptok/gocomm/common"
	// "github.com/tiptok/gocomm/pkg/cache/model"
	// "github.com/tiptok/gocomm/pkg/log"
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/linmadan/egglib-go/persistent/cache"
	"github.com/linmadan/egglib-go/persistent/cache/model"
	"github.com/linmadan/egglib-go/web/beego/utils"
)

const (
	apqExtension  = "apq"
	defaultExpire = 60 // 60s
)

// AtomicPersistenceHandler  if routers match , atomic persistence response data to cache store,cache will be used in future lookups
// links:https://gqlgen.com/reference/apq/
// links:https://github.com/apollographql/apollo-link-persisted-queries
func AtomicPersistenceQueryHandler(options ...option) func(next web.FilterFunc) web.FilterFunc {
	option := NewOptions(options...)
	option.ValidAPQ()
	return func(next web.FilterFunc) web.FilterFunc {
		return func(ctx *context.Context) {
			var (
				queryHash string
				err       error
			)
			r := ctx.Request
			w := ctx.ResponseWriter
			if !(r.Method == http.MethodPost || r.Method == http.MethodGet) {
				next(ctx)
				return
			}
			if !checkRouter(r, option.routers) {
				next(ctx)
				return
			}
			if queryHash, err = ComputeHttpRequestQueryHash(r); err != nil {
				fmt.Println(err)
				next(ctx)
				return
			}
			var item string
			// if cache is miss , store the newest data to cache
			if v, err := option.cache.Get(redisKey(option.serviceName, queryHash), &item); err != nil || v == nil {
				if err != nil {
					fmt.Println(err)
				}
				responseBuf := bytes.NewBuffer(nil)
				crw := newCacheResponseWrite(w, responseBuf)
				ctx.ResponseWriter = &context.Response{
					ResponseWriter: crw,
					Started:        ctx.ResponseWriter.Started,
					Status:         ctx.ResponseWriter.Status,
					Elapsed:        ctx.ResponseWriter.Elapsed,
				}
				next(ctx)
				ctx.ResponseWriter = crw.writer
				if err := option.cache.Set(redisKey(option.serviceName, queryHash), model.NewItem(responseBuf.String(), option.expire)); err != nil {
					fmt.Println(err)
				}
				return
			}
			// 此处不能提前设置状态，否则beego内部框架会识别response已被处理,导致content-type:text-plain(一直是)
			// 详见 :https://blog.csdn.net/yes169yes123/article/details/103126655
			// w.WriteHeader(http.StatusAccepted)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(item))
		}
	}
}

// computeQueryHash compute hash key
func ComputeQueryHash(query string) string {
	b := sha256.Sum256([]byte(query))
	return hex.EncodeToString(b[:])
}

//ComputeHttpRequestQueryHash  compute request query hash
func ComputeHttpRequestQueryHash(r *http.Request) (string, error) {
	var queryHash string
	if r.Method == http.MethodGet {
		queryHash = ComputeQueryHash(r.URL.String())
	} else if r.Method == http.MethodPost {
		var body []byte
		var err error
		body, r.Body , err = utils.DumpReadCloser(r.Body)
		if err != nil {
			return "", err
		}
		queryHash = ComputeQueryHash(r.URL.String() + string(body))
	}
	return queryHash, nil
}

func checkRouter(r *http.Request, routers []string) bool {
	for i := range routers {
		if utils.KeyMatch3(r.URL.Path, routers[i]) {
			return true
		}
	}
	return false
}

func redisKey(serviceName, hash string) string {
	return strings.Join([]string{serviceName, apqExtension, hash}, ":")
}

// cacheResponseWrite buffer response data in future use
type cacheResponseWrite struct {
	writer *context.Response
	buf    *bytes.Buffer
}

func (w *cacheResponseWrite) Header() http.Header {
	return w.writer.Header()
}

func (w *cacheResponseWrite) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *cacheResponseWrite) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}

func newCacheResponseWrite(writer *context.Response, buf *bytes.Buffer) *cacheResponseWrite {
	return &cacheResponseWrite{
		writer: writer,
		buf:    buf,
	}
}

type Options struct {
	hashFunc         func(query string) string
	requestQueryHash func(r *http.Request) (string, error)
	serviceName      string
	cache            cache.Cache
	expire           int
	routers          []string
}

type option func(options *Options)

func WithRequestQueryHashFunc(requestQueryHash func(r *http.Request) (string, error)) option {
	return func(options *Options) {
		options.requestQueryHash = requestQueryHash
	}
}

func WithServiceName(serviceName string) option {
	return func(options *Options) {
		options.serviceName = serviceName
	}
}

func WithCache(cache cache.Cache) option {
	return func(options *Options) {
		options.cache = cache
	}
}

// WithExpire set cache expire duration (unit:second)
func WithExpire(expire int) option {
	return func(options *Options) {
		options.expire = expire
	}
}

func WithRouters(routers []string) option {
	return func(options *Options) {
		options.routers = routers
	}
}

func NewOptions(options ...option) *Options {
	option := &Options{
		hashFunc:         ComputeQueryHash,
		requestQueryHash: ComputeHttpRequestQueryHash,
		expire:           defaultExpire,
	}
	for i := range options {
		options[i](option)
	}
	return option
}

func (o *Options) ValidAPQ() error {
	if o.cache == nil {
		panic("Options cache is null")
	}
	if len(o.serviceName) == 0 {
		panic("Options serviceName is empty")
	}
	return nil
}
