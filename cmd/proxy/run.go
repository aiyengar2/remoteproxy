package proxy

import (
	"context"

	"github.com/aiyengar2/portexporter/pkg/config"
	"github.com/aiyengar2/portexporter/pkg/proxy"
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
			Value: ":8080",
		},
		cli.StringFlag{
			Name:  "cert-file",
			Usage: "A file containing a TLS cert used to set up TLS encrypted proxy connections",
		},
		cli.StringFlag{
			Name:  "key-file",
			Usage: "A file containing a TLS key used to set up TLS encrypted proxy connections",
		},
		cli.StringFlag{
			Name:  "cacert-file",
			Usage: "A file containing a caCert to be used to verify incoming TLS encrypted proxy connections",
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
	listen := cliCtx.String("listen")
	certFile := cliCtx.String("cert-file")
	keyFile := cliCtx.String("key-file")
	caCertFile := cliCtx.String("cacert-file")
	debug := cliCtx.Bool("debug")
	printTunnelData := cliCtx.Bool("print-tunnel-data")

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		remotedialer.PrintTunnelData = printTunnelData
	}

	cfg := config.TLSServer{
		CertFile:   certFile,
		KeyFile:    keyFile,
		CaCertFile: caCertFile,
	}
	s := proxy.NewServer(listen, cfg)

	return s.Start(ctx)
}
