package lotus

import (
	"net/http"

	"github.com/go-chi/chi"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/types"
)

// Router handler for filecoin deal
func Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/deal", func(r chi.Router) {
		r.Post("/push", apicontext.Bind(makeDeal, types.MakeDealOptions{})) // 发布交易
		r.Get("/clean", apicontext.Bind(cleanDeal))                         // 手动触发同步交易状态
		r.Get("/status", apicontext.Bind(status))                           // 手动触发同步交易状态
	})
	return r
}

func status(ctx *apicontext.APIContext) {
	// pool := lotus.GetDealPool(ctx)
	// ctx.JSON(fmt.Sprintf("%d\n", pool.Queue.Length()))
	ctx.JSON("success")
}

func cleanDeal(ctx *apicontext.APIContext) {
	// pool := lotus.GetDealPool(ctx)
	// pool.CleanQueue()
	ctx.JSON("success")
}

func makeDeal(ctx *apicontext.APIContext, ops *types.MakeDealOptions) {
	// pool := lotus.GetDealPool(ctx)
	// if err := pool.AddTask(ops); err != nil {
	// ctx.Error(render.ServerError, err, "addTask")
	// return
	// }
	ctx.JSON("success")
}
