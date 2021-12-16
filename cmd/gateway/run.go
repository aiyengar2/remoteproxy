package gateway

import (
	"context"

	"github.com/aiyengar2/portexporter/pkg/gateway"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/rancher/remotedialer"
)

var (
	runFlags = []cli.Flag{
		cli.StringFlag{
			Name:     "proxy-url",
			Usage:    "The address of the proxy that the gateway must make an outbound connection to (e.g. wss://port-exporter-proxy.port-exporter.svc.cluster.local:8080/connect)",
			Required: true,
		},
		cli.StringFlag{
			Name:  "cacert-file",
			Usage: "A file containing a TLS cacert used to verify the TLS certs provided by the proxy when setting up a TLS encrypted proxy connection",
		},
		cli.BoolFlag{
			Name:  "insecure-skip-verify",
			Usage: "Whethert to skip verifying certs provided by the proxy when setting up a TLS encrypted proxy connection",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
		cli.BoolFlag{
			Name:  "print-tunnel-data",
			Usage: "Enable printing remotedialer tunnel data. Requires debug to be set to true",
		},
	}
)

func run(cliCtx *cli.Context) (err error) {
	ctx := signals.SetupSignalHandler(context.Background())

	// parse flags
	proxyUrl := cliCtx.String("proxy-url")
	expose := cliCtx.StringSlice("expose")
	caCertFile := cliCtx.String("cacert-file")
	insecureSkipVerify := cliCtx.Bool("insecure-skip-verify")
	debug := cliCtx.Bool("debug")
	printTunnelData := cliCtx.Bool("print-tunnel-data")

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		remotedialer.PrintTunnelData = printTunnelData
	}

	cfg := gateway.Config{
		Expose: expose,
	}

	cfg.InsecureSkipVerify = insecureSkipVerify
	cfg.CaCertFile = caCertFile

	g := gateway.NewServer(proxyUrl, cfg)

	return g.Start(ctx)
}
