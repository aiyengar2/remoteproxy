package proxy

import (
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "proxy",
		Usage:  "Proxies incoming HTTP(s) requests to registered Gateways",
		Action: run,
		Flags:  runFlags,
	}
}
