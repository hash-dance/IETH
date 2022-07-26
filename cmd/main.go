package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/guowenshuai/ieth/cmd/export"
	"github.com/guowenshuai/ieth/cmd/ipfs"
	"github.com/guowenshuai/ieth/cmd/lotus"
	"github.com/urfave/cli/v2"
)

var (
	// VERSION current version
	VERSION = "v0.1.0"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd := []*cli.Command{
		ipfs.IpfsCmd,
		lotus.LotusCmd,
		export.ExportCmd,
	}

	app := &cli.App{
		Name:    "ieth-cli",
		Version: VERSION,
		Usage:   "everything to ipfs and filecoin!",
		Flags: []cli.Flag{
			// 配置认证目录
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"IETH_PATH"},
				Hidden:  true,
				Value:   "~/.ieth", //
			},
		},
		// 配置命令行参数
		Commands: cmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
}
