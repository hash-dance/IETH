package report

import (
	"net/http"

	"github.com/go-chi/chi"
	apicontext "github.com/guowenshuai/ieth/modules/context"
)

// Router handler for export
func Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/report", apicontext.Bind(export))
	return r
}

func export(ctx *apicontext.APIContext) {
	// all := param.QueryString(ctx.Req, "all")
	// checkAll := false
	// if all != "" {
	// 	checkAll = true
	// }
	// if  allDeals, err := report.Export(ctx, checkAll); err != nil {
	// 	ctx.Error(render.ServerError, err, "Export")
	// 	return
	// } else {
	// 	ctx.JSON(allDeals)
	// }
}
