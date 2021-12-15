package gateway

import (
	"context"

	"net/http"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/urfave/cli"

	"github.com/rancher/remotedialer"
)

var (
	runFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "proxy-url",
			Usage: "A URL pointing to the proxy (e.g. wss://port-exporter-proxy.port-exporter.svc.cluster.local:8080/connect)",
		},
		cli.StringSliceFlag{
			Name:  "addresses",
			Usage: "Comma-delimited list of addresses to expose via the gateway",
		},
	}
)

type Mapping struct {
	SrcAddr string
	DstAddr string
}

func run(cliCtx *cli.Context) (err error) {
	ctx := signals.SetupSignalHandler(context.Background())

	proxyUrl := cliCtx.String("proxy-url")

	addresses := cliCtx.StringSlice("addresses")

	addressMap := make(map[string]bool)
	for _, address := range addresses {
		addressMap[address] = true
	}

	auth := func(proto, address string) bool {
		if proto != "tcp" {
			return false
		}
		return addressMap[address]
	}

	onConnect := func(ctx context.Context, _ *remotedialer.Session) error {
		return nil
	}

	return remotedialer.ClientConnect(ctx, proxyUrl, http.Header{}, nil, auth, onConnect)
}
