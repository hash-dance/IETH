package lotus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/guowenshuai/ieth/cmd/client"
	"github.com/urfave/cli/v2"
)

var cleanCmd = &cli.Command{
	Name:    "clean",
	Usage:   "clean deal pool",
	Action: func(c *cli.Context) error {
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("获取客户端失败  %s\n", err.Error())
			return err
		}
		ret, res, err := cli.LotusApi.LotusDealCleanGet(context.Background(), nil)
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

var statusCmd = &cli.Command{
	Name:    "status",
	Usage:   "status for deal pool",
	Action: func(c *cli.Context) error {
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("获取客户端失败  %s\n", err.Error())
			return err
		}
		ret, res, err := cli.LotusApi.LotusDealStatusGet(context.Background(), nil)
		if err != nil {
			fmt.Printf("LotusDealCleanGet get err %s", err.Error())
			return err
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("请求失败")
			return err
		}
		fmt.Printf("%s", ret.Data)
		return nil
	},
}
