package lotus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antihax/optional"
	"github.com/guowenshuai/ieth/cmd/client"
	"github.com/urfave/cli/v2"
)

// 存储几份
// Duplicate int32 `json:"duplicate"`
// // 存储时长
// Duration int32 `json:"duration"`
// // 存储单价
// Price string `json:"price"`
// // 发单到存储市场
// Market bool `json:"market"`
// // 钱包地址
// Wallet string `json:"wallet"`
// // 发送多少单
// Nums int32 `json:"nums"`

var pushCmd = &cli.Command{
	Name:  "make",
	Usage: "make some deal",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "duration",
			Usage: "存储时长(天)",
			Value: 180,
		},
		&cli.StringFlag{
			Name:  "price",
			Usage: "存储单价 FIL",
			Value: "0",
		},
		&cli.BoolFlag{
			Name:    "market",
			Aliases: nil,
			Usage:   "发布到存储市场",
		},
		&cli.BoolFlag{
			Name:    "sort",
			Aliases: nil,
			Usage:   "默认按文件大小降序发单",
		},
		&cli.StringFlag{
			Name:     "wallet",
			Usage:    "钱包地址",
			Value:    "",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "nums",
			Usage: "发送多少单",
			Value: 100,
		},
	},
	Action: func(c *cli.Context) error {
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("获取客户端失败  %s\n", err.Error())
			return err

		}
		ret, res, err := cli.LotusApi.LotusDealPushPost(context.Background(), &client.LotusApiLotusDealPushPostOpts{
			Root: optional.NewInterface(client.EmptyObject3{
				Duplicate: 0,
				Duration:  int32(c.Int("duration")),
				Price:     c.String("price"),
				Market:    c.Bool("market"),
				Wallet:    c.String("wallet"),
				Nums:      int32(c.Int("nums")),
				Sort:      c.Bool("sort"),
			}),
		})
		if err != nil {
			fmt.Printf("LotusDealCleanGet get err %s", err.Error())
			return err
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("请求失败")
			return err
		}
		fmt.Printf("%s\n", ret.Data)
		return nil
	},
}
