package lotus

import (
	"github.com/urfave/cli/v2"
)

var LotusCmd = &cli.Command{
	Name:  "lotus",
	Usage: "Tools for deal with lotus miner",
	Subcommands: []*cli.Command{
		deal,
		priceCmd,
	},
	Flags: []cli.Flag{},
}

var deal = &cli.Command{
	Name: "deal",
	Usage: "make deal with lotus miner",
	Subcommands: []*cli.Command{
		pushCmd,
		cleanCmd,
		statusCmd,
	},
}