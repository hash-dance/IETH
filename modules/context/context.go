// Package apicontext defined common context
package apicontext

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/filecoin-project/lotus/api"
	"github.com/guowenshuai/ieth/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/guowenshuai/ieth/modules/common/render"
)

// APIContext struct
type APIContext struct {
	Context     context.Context
	MongoClient *mongo.Database
	Config      *types.Config
	FullNode    api.FullNode
	Req         *http.Request
	Writer      http.ResponseWriter
}

func (ctx *APIContext) Redirect(url string, code int) {
	http.Redirect(ctx.Writer, ctx.Req, url, code)
}

func (ctx *APIContext) JSON(data interface{}) {
	render.SendJSON(ctx.Writer, ctx.Req, data)
}

func (ctx *APIContext) JSONPagination(data interface{}) {
	render.SendPaginationJSON(ctx.Writer, ctx.Req, data)
}

func (ctx *APIContext) Error(code render.ErrorCode, err error, message string) {
	e := errors.Wrap(err, message)
	logrus.Error(e.Error())
	render.SendError(ctx.Writer, ctx.Req, code, e)
}

func (ctx *APIContext) Errorf(code render.ErrorCode, err error, format string, args ...interface{}) {
	e := errors.Wrapf(err, format, args)
	logrus.Error(e.Error())
	render.SendError(ctx.Writer, ctx.Req, code, e)
}

type contextKey struct {
	name string
}

// Middleware load common context
func Middleware(apiContext *APIContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKey{"api-context"}, &APIContext{
				Context:     apiContext.Context,
				MongoClient: apiContext.MongoClient,
				Config:      apiContext.Config,
				FullNode:    apiContext.FullNode,
				Req:         nil,
				Writer:      nil,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// ReadAPIContext load apicontext from context
func ReadAPIContext(ctx context.Context) *APIContext {
	return ctx.Value(contextKey{"api-context"}).(*APIContext)
}

// Bind read input and valida input fields
func Bind(handler interface{}, input ...interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := ReadAPIContext(r.Context())
		ctx.Writer = w
		ctx.Req = r
		defer func(ctx *APIContext) {
			if err := recover(); err != nil {
				ctx.Error(render.ServerError, fmt.Errorf("bind error: [%s]", err), "defer")
			}
		}(ctx)

		contentType := ctx.Req.Header.Get("Content-Type")
		if ctx.Req.Method == "POST" || ctx.Req.Method == "PUT" || len(contentType) > 0 {
			if len(input) == 0 { // no request body input
				fn := reflect.ValueOf(handler)
				fn.Call([]reflect.Value{reflect.ValueOf(ctx)})
				// ctx.Error(render.InvalidData, fmt.Errorf("must body"), "apiContext bind")
				return
			}
			obj := input[0]
			typ := reflect.TypeOf(obj)
			data := reflect.New(typ).Interface()
			// todo parse contentType
			// 		switch {
			// 		case strings.Contains(contentType, "form-urlencoded"):
			// 		case strings.Contains(contentType, "multipart/form-data"):
			// 		case strings.Contains(contentType, "json"):
			// 		default:
			// 			err := fmt.Errorf("bind parse error")
			// 			logrus.Error(err.Error())
			// 		}
			if err := render.DecodeJSON(r.Body, data); err != nil {
				ctx.Error(render.InvalidData, err, "decode json")
				return
			}
			fn := reflect.ValueOf(handler)
			fn.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(data)})
		} else { // get request
			fn := reflect.ValueOf(handler)
			fn.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
	}
}
