package proxy

import (
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:    "proxy",
		Aliases: []string{"proxy"},
		Usage:   "Creates a proxy that listens for registering gateways and proxies incoming requests to the corresponding gateways",
		Action:  run,
		Flags:   runFlags,
	}
}
