package ipfs

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/guowenshuai/ieth/modules/common/render"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/modules/ipfs"
	"github.com/guowenshuai/ieth/types"
)

// Router handler for ipfs
func Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/push", apicontext.Bind(push, types.IPFSPushOptions{}))
	r.Post("/list", apicontext.Bind(list, types.IPFSListOptions{}))
	return r
}

func push(ctx *apicontext.APIContext, opt *types.IPFSPushOptions) {
	go ipfs.Push(ctx, opt.Path, opt.Recursive)
	ctx.JSON("success")
}

func list(ctx *apicontext.APIContext, opt *types.IPFSListOptions) {
	links, err :=  ipfs.List(ctx, opt.Cid, opt.Recursive)
	if err != nil {
		ctx.Error(render.ServerError, err, "ipfs list links")
		return
	}
	ctx.JSON(links)
}
