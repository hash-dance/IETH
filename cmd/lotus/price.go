package lotus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/guowenshuai/ieth/cmd/client"
	"github.com/urfave/cli/v2"
)

var priceCmd = &cli.Command{
	Name:  "price",
	Usage: "Tools for miner price",
	// Flags: []cli.Flag{
	// 	&cli.BoolFlag{
	// 		Name:    "limit",
	// 		Aliases: []string{"l"},
	// 		Usage:   "limit lines",
	// 		Value:   false,
	// 	},
	// },
	Action: func(c *cli.Context) error {
		cli, err := client.NewCli(c)
		if err != nil {
			fmt.Printf("失败 %s\n", err.Error())
			return err
		}
		ret, res, err := cli.PriceApi.PriceGet(context.Background(), nil)
		if err != nil {
			fmt.Printf("price get err %s", err.Error())
			return err
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("请求失败")
			return err
		}
		fmt.Printf("miner Price_per_GiB Verified_Price_per_GiB Min_size Max_size Timestamp Expiry SeqNo\n")

		for _, ask :=  range ret.Data {
			fmt.Printf("%s %s %s %s %s %s %s %s\n", ask.Miner, ask.PricePerGib, ask.VerifiedPricePerGib, ask.MinSize,
				ask.MaxSize, ask.Timestamp, ask.Expiry, ask.Seqno)
		}
		return nil
	},
}
