package ipfs

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/antihax/optional"
	"github.com/guowenshuai/ieth/cmd/client"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

var pushCmd = &cli.Command{
	Name:    "push",
	Aliases: []string{"a"},
	Usage:   "push a file or dir to ipfs",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "push recursive for dirs, not entire dir",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		if  c.Args().First() == "" {
			return fmt.Errorf("必须指定路径\n")
		}
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("获取客户端失败 %s\n", err.Error())
			return err
		}

		p, err := homedir.Expand(c.Args().First())
		if err != nil {
			fmt.Printf("解析路径失败 %s\n", err.Error())
			return nil
		}
		p1, err := filepath.Abs(p)
		if err != nil {
			fmt.Printf("解析路径失败 %s\n", err.Error())
			return nil
		}

		ret, res, err := cli.IpfsApi.IpfsPushPost(context.Background(), &client.IpfsApiIpfsPushPostOpts{
			Root: optional.NewInterface(client.EmptyObject1{
				Recursive: c.Bool("recursive"),
				Path:      p1,
			}),
		})
		if err != nil {
			fmt.Printf("IpfsPushPost get err %s", err.Error())
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
