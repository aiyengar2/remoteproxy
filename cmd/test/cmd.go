package test

import (
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "test",
		Usage:  "Spins up a simple hello-world server for testing at a given port",
		Action: run,
		Flags:  runFlags,
	}
}
