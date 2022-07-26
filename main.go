// Package main ieth
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
package main

import (
	"context"
	"os"
	"runtime"

	"github.com/guowenshuai/ieth/conf"
	db "github.com/guowenshuai/ieth/db/mongo"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	lotus2 "github.com/guowenshuai/ieth/modules/lotus"
	"github.com/guowenshuai/ieth/modules/util"
	"github.com/guowenshuai/ieth/routers"
	"github.com/guowenshuai/ieth/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	// VERSION current version
	VERSION = "v0.1.0"
)

func initConfig() *types.Config {
	// 初始化启动配置
	config := &types.Config{
		Mongodb: &types.Mongodb{
			Server:   "192.168.1.157:27037",
			NoAuth:   false,
			Username: "admin",
			Password: "admin",
			Database: "ieth",
		},
		Lotus: &types.Lotus{
			Token:   "",
			Address: "",
		},
		Ipfs: &types.Ipfs{
			Token:   "",
			Address: "",
		},
		BaseConf: &types.BaseConf{
			Repo:           "~/.ieth",
			Debug:          true,
			Cors:           true,
			LogFormat:      conf.LogsText,
			LogPath:        "./log",
			LogDispatch:    false,
			HTTPListenPort: 80,
			Monitor:        false,
			SessionTimeOut: 30,
			SSL:            false,
			SSLCrtFile:     "",
			SSLKeyFile:     "",
		},
		Setting: &types.Setting{
			DealTimeout:      30,
			MaxDealTransfers: 10,
			MinerMaxDeals:    2,
			MaxDealOne:       10,
			Wallet:           "",
			Duration:         366,
		},
	}
	if err := conf.LoadYaml("conf.yaml", config); err != nil {
		panic(err)
	}
	return config
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := &cli.App{
		Name:    "ieth",
		Version: VERSION,
		Usage:   "everything to ipfs and filecoin!",
	}

	ctx := util.SigTermCancelContext(context.Background())
	config := initConfig()
	nodeApi, closer, err := lotus2.NewLocalFullNodeRPC(config)
	defer closer()
	if err != nil {
		logrus.Fatalf("connecting with lotus failed: %s\n", err)
	}
	client, err := db.Connect(ctx, config)
	if err != nil {
		logrus.Fatalf("connecting mongo err: %s\n", err)
	}
	conf.InitLogs(config) // 初始化日志配置

	// 构建apicontext
	apiContext := apicontext.APIContext{
		Context:     ctx,
		Config:      config,
		FullNode:    nodeApi,
		MongoClient: client,
	}

	app.Action = func(c *cli.Context) error {
		return run(&apiContext)
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx *apicontext.APIContext) error {
	server := routers.NewServerConfig(ctx)
	// watcher.Watch(ctx)
	server.Build()
	// lotus askPrice 服务
	// go lotus2.AskPrice(ctx)
	go lotus2.SyncDeals(ctx)
	// if err := lotus2.StartDealPool(ctx); err != nil {
	// 	logrus.Fatalf("start deal pool err %s", err.Error())
	// }
	<-ctx.Context.Done()
	db.Disconnect()
	return nil
}
