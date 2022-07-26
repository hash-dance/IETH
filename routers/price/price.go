package price

import (
	"net/http"

	"github.com/go-chi/chi"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/modules/lotus"
)

// Router handler for price
func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", apicontext.Bind(list))
	return r
}

func list(ctx *apicontext.APIContext) {
	ctx.JSON(lotus.GetMarketStorageAsk())
}
