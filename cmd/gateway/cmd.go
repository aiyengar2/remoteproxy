package gateway

import (
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "gateway",
		Usage:  "Registers with a Proxy and proxies / reverse proxies HTTP(s) connections on the local network",
		Action: run,
		Flags:  runFlags,
	}
}
