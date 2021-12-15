package gateway

import (
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "gateway",
		Usage:  "Registers with a proxy and opens a gateway into your network for incoming proxied requests",
		Action: run,
		Flags:  runFlags,
	}
}
