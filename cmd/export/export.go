package export

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antihax/optional"
	"github.com/guowenshuai/ieth/cmd/client"
	"github.com/urfave/cli/v2"
)

var ExportCmd = &cli.Command{
	Name:  "export",
	Usage: "Tools for report",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "export all onChain deals",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("失败 %s\n", err.Error())
			return err
		}
		ops := client.ExportApiExportReportGetOpts{}
		if c.Bool("all") {
			ops.All = optional.NewString("true")
		}
		ret, res, err := cli.ExportApi.ExportReportGet(context.Background(), &ops)
		if err != nil {
			fmt.Printf("price get err %s", err.Error())
			return err
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("请求失败")
			return err
		}
		for _, ask := range ret.Data {
			fmt.Println(ask)
		}
		return nil
	},
}
