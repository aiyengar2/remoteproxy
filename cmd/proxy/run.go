package proxy

import (
	"context"
	"fmt"
	"os"

	"github.com/aiyengar2/portexporter/pkg/defaults"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	runFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "server-url",
			Usage: "The address of the server that the client must make an outbound connection to (e.g. ws://localhost:10123)",
			Value: "ws://localhost:10123",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "[optional] Specifies the path of the configuration",
			Value: defaults.ConfigPath,
		},
	}

	c *client
)

func parseFlags(cliCtx *cli.Context) (err error) {
	// Get hostname to identify backend connection
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get hostname: %s", err)
	}
	proxyURL := cliCtx.String("server-url")
	listenAndReverseProxy := map[port]*reverseProxy{
		port(8080): nil,
	}

	c = &client{
		clientKey:             hostname,
		proxyURL:              proxyURL,
		listenAndReverseProxy: listenAndReverseProxy,
	}

	return nil
}

func run(cliCtx *cli.Context) (err error) {
	logrus.Infof("Initializing client...")

	ctx := signals.SetupSignalHandler(context.Background())
	if err := c.Start(ctx); err != nil {
		return err
	}

	return nil
}
