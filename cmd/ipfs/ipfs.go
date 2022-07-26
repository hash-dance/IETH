package ipfs

import (
	"github.com/urfave/cli/v2"
)

var IpfsCmd = &cli.Command{
	Name:  "ipfs",
	Usage: "Tools for deal with ipfs node",
	Subcommands: []*cli.Command{
		pushCmd,
		listCmd,
	},
	Flags: []cli.Flag{},
}
