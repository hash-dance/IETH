/*Package server defined server config and run server
 */
package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof" // pprof
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/guowenshuai/ieth/modules/repo"
	"github.com/guowenshuai/ieth/types"
	_ "github.com/mkevac/debugcharts" // debugcharts
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/guowenshuai/ieth/conf"
	apicontext "github.com/guowenshuai/ieth/modules/context"
)

// Config Defined struct of router
type Config struct {
	Context *apicontext.APIContext
	Timeout time.Duration

	AuthRouter   map[string]http.Handler
	PublicRouter map[string]http.Handler
	CustomRouter map[string]http.Handler
}

const (
	// APIVersion api version
	APIVersion = "/v1"
	// APIPublicVersion public version
	APIPublicVersion = "/v1-public"
)

var (
	debugIndex = `<html>
				<head><title>raging-server monitor</title></head>
				<body>
				   <h1>raging server monitor</h1>
				   <p><a href='/metrics'>metrics</a></p>
				   <p><a href='/debug/pprof'>pprof</a></p>
				   <p><a href='/debug/charts'>debugcharts</a></p>
				   </body>
				</html>
			  `
)

// Build init server config and run
func (c *Config) Build() {
	r := chi.NewRouter()
	config := c.Context.Config
	configMiddleware(c, r)

	r.Group(func(r chi.Router) {
		// first-login middleware
		r.Group(func(r chi.Router) {
			// 认证中间件
			r.Mount(APIVersion, registryRouters(c.AuthRouter))
			r.Mount("/", registryRouters(c.CustomRouter))
		})
		r.Mount(APIPublicVersion, registryRouters(c.PublicRouter))

		if config.BaseConf.Monitor {
			http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(debugIndex))
				if err != nil {
					return
				}
			})
			http.Handle("/metrics", promhttp.Handler())
			r.Mount("/debug/pprof", http.DefaultServeMux)
			r.Mount("/debug", http.DefaultServeMux)
			r.Mount("/metrics", http.DefaultServeMux)
		}
	})
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
	start(config, r)
}

func registryRouters(routers map[string]http.Handler) *chi.Mux {
	r := chi.NewRouter()
	for pattern, router := range routers {
		r.Mount(pattern, router)
	}
	return r
}

func configMiddleware(c *Config, r *chi.Mux) {
	r.Use(middleware.RealIP, middleware.Recoverer, middleware.Timeout(c.Timeout))

	config := c.Context.Config
	r.Use(middleware.RequestID, conf.RequestLogger()) // 生成requestID RequestIDKey, 插入日志中间件

	if config.BaseConf.Cors {
		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		cors := cors.New(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		r.Use(cors.Handler)
	}

	// apicontext初始化
	r.Use(apicontext.Middleware(c.Context))
}

// run server
func start(config *types.Config, mux http.Handler) {

	rp, err := repo.NewRepo(config.BaseConf.Repo)
	if err != nil {
		logrus.Fatal(err)
	}

	if config.BaseConf.SSL {
		cf, err := ioutil.ReadFile(config.BaseConf.SSLCrtFile)
		if err != nil {
			panic(err)
		}
		key, err := ioutil.ReadFile(config.BaseConf.SSLKeyFile)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair(cf, key)
		if err != nil {
			panic("can not create tls client: " + err.Error())
		}
		server := http.Server{
			Addr:    ":" + strconv.Itoa(config.BaseConf.HTTPListenPort),
			Handler: mux,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
		// 启动http服务
		go func() {
			logrus.Fatal(server.ListenAndServeTLS("", ""))
		}()
		if err := rp.SetApi(fmt.Sprintf("https://127.0.0.1:%d/v1", config.BaseConf.HTTPListenPort)); err != nil {
			logrus.Fatal(err)
		}
	} else {
		// 启动http服务
		go func() {
			logrus.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.BaseConf.HTTPListenPort), mux))
		}()
		if err := rp.SetApi(fmt.Sprintf("http://127.0.0.1:%d/v1", config.BaseConf.HTTPListenPort)); err != nil {
			logrus.Fatal(err)
		}
	}
}