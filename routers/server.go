/*Package routers registry all routers
 */
package routers

import (
	"net/http"
	"time"

	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/modules/server"
	"github.com/guowenshuai/ieth/routers/ipfs"
	"github.com/guowenshuai/ieth/routers/lotus"
	"github.com/guowenshuai/ieth/routers/price"
	"github.com/guowenshuai/ieth/routers/report"
)

// NewServerConfig return a server config
func NewServerConfig(ctx *apicontext.APIContext) *server.Config {
	return &server.Config{
		Context: ctx,
		Timeout: 60 * time.Second,
		AuthRouter: map[string]http.Handler{
			"/price":  price.Router(),
			"/ipfs":   ipfs.Router(),
			"/lotus":  lotus.Router(),
			"/export": report.Router(),
		},
		PublicRouter: map[string]http.Handler{},
		CustomRouter: map[string]http.Handler{
			// "/": host.Router(),
		},
	}
}
