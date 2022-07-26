package ipfs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antihax/optional"
	"github.com/guowenshuai/ieth/cmd/client"
	"github.com/urfave/cli/v2"
)

var listCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List links from an cid",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "list recursive for dirs",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		if  c.Args().First() == "" {
			return fmt.Errorf("必须指定cid\n")
		}
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("获取客户端失败  %s\n", err.Error())
			return err
		}
		ret, res, err := cli.IpfsApi.IpfsListPost(context.Background(), &client.IpfsApiIpfsListPostOpts{
			Root: optional.NewInterface(client.EmptyObject{
				Recursive: c.Bool("recursive"),
				Cid:       c.Args().First(),
			}),
		})
		if err != nil {
			fmt.Printf("IpfsListPost get err %s", err.Error())
			return err
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("请求失败")
			return err
		}
		fmt.Printf("name hash size type\n")

		for _, fil :=  range ret.Data {
			fmt.Printf("%s %s %d %d\n", fil.Name, fil.Hash, int64(fil.Size), int64(fil.Type_))
		}
		return nil
	},
}
