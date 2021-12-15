package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/aiyengar2/portexporter/pkg/defaults"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
	"inet.af/tcpproxy"
)

type client struct {
	clientKey             string
	proxyURL              string
	listenAndReverseProxy map[port]*reverseProxy
}

func (c *client) Start(ctx context.Context) error {
	headers := http.Header{
		defaults.ClientKeyHeader: []string{c.clientKey},
	}

	addressAllowed := make(map[address]bool, len(c.listenAndReverseProxy))
	addressForward := make(map[address]address, len(c.listenAndReverseProxy))

	for p, rp := range c.listenAndReverseProxy {
		listenAddress := allHosts(p)
		addressAllowed[listenAddress] = true

		if rp == nil {
			// default behavior is to forward to localhost:p
			addressForward[listenAddress] = localhost(p)
			continue
		}

		// start the reverse proxy and forward all requests to :rp.ListenTo
		logrus.Infof("Registering %s", rp)
		addressForward[listenAddress] = allHosts(rp.listenTo)
		rpCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		go rp.Start(rpCtx)
	}

	connAuth := func(proto, addr string) bool {
		if proto != "tcp" {
			return false
		}
		if exists, ok := addressAllowed[address(addr)]; exists && ok {
			return true
		}
		return false
	}

	onConn := func(ctx context.Context, s *remotedialer.Session) error {
		proxy := &tcpproxy.Proxy{}
		for listenAddress, forwardAddress := range addressForward {
			proxy.AddRoute(string(listenAddress), &tcpproxy.DialProxy{
				DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
					return s.Dial(ctx, "tcp", string(forwardAddress))
				},
			})
		}
		if err := proxy.Start(); err != nil {
			return err
		}
		<-ctx.Done()
		proxy.Close()
		return nil
	}

	return remotedialer.ConnectToProxy(ctx, c.proxyURL, headers, connAuth, nil, onConn)
}

type reverseProxy struct {
	listenTo  port
	forwardTo address
}

func (r *reverseProxy) Start(ctx context.Context) {
	defer r.restartIfClosed(ctx)

	logrus.Infof("%s Sleeping for 5 seconds...", r)
	time.Sleep(time.Second * 5)
}

func (r *reverseProxy) restartIfClosed(ctx context.Context) {
	logrus.Warnf("%s Closed, attempting to restart...", r)
	select {
	case <-ctx.Done():
		return
	default:
		r.Start(ctx)
	}
}

func (r reverseProxy) String() string {
	return fmt.Sprintf("ReverseProxy[:%d => %s]", r.listenTo, r.forwardTo)
}
