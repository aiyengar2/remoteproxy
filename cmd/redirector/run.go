package redirector

import (
	"context"

	"github.com/aiyengar2/portexporter/pkg/redirect"
	"github.com/rancher/remotedialer"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	runFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Usage: "The address to listen to incoming HTTP requests on",
			Value: ":8081",
		},
		cli.StringFlag{
			Name:      "config",
			Usage:     "The location of a configuration file for the redirector (default: redirect.yaml)",
			TakesFile: true,
			Value:     redirect.DefaultRedirectConfigFile,
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
	}
)

func run(cliCtx *cli.Context) (err error) {
	ctx := signals.SetupSignalHandler(context.Background())

	// parse flags
	listen := cliCtx.String("listen")
	config := cliCtx.String("config")
	debug := cliCtx.Bool("debug")

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		remotedialer.PrintTunnelData = true
	}

	cfg, err := redirect.Load(config)
	if err != nil {
		logrus.Fatal(err)
	}
	s := redirect.NewServer(listen, cfg)

	return s.Start(ctx)
}
