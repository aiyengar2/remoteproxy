package redirector

import (
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "redirector",
		Usage:  "Redirects the contents of HTTP(s) requests to a provided address",
		Action: run,
		Flags:  runFlags,
	}
}
