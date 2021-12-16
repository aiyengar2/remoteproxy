package gateway

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/aiyengar2/portexporter/pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
)

type gatewayServer struct {
	proxyUrl  string
	expose    []string
	tlsConfig *tls.Config
}

func NewServer(proxyUrl string, config Config) *gatewayServer {
	s := &gatewayServer{
		proxyUrl: proxyUrl,
		expose:   config.Expose,
	}
	if strings.HasPrefix(s.proxyUrl, "wss://") {
		s.tlsConfig = config.TLSConfig(proxyUrl)
	}
	return s
}

func (s *gatewayServer) Start(ctx context.Context) error {
	ip := utils.GetHostIP()
	logrus.Infof("Using id [%s]", ip)

	headers := http.Header{
		"X-Proxy-Tunnel-ID": []string{ip},
	}
	connAuth := getConnectAuthorizer(s.expose)
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: remotedialer.HandshakeTimeOut,
		TLSClientConfig:  s.tlsConfig,
	}
	return remotedialer.ClientConnect(ctx, s.proxyUrl, headers, dialer, connAuth, onConnect)
}

func getConnectAuthorizer(expose []string) remotedialer.ConnectAuthorizer {
	var addressMap map[string]bool
	if len(expose) > 0 {
		addressMap = make(map[string]bool)
		for _, address := range expose {
			addressMap[address] = true
		}
	}

	return func(proto, address string) bool {
		logrus.Debugf("Received request to %s://%s", proto, address)
		if proto != "tcp" {
			// only tcp is supported
			return false
		}
		// if addressMap is nil, then we expose everything by default
		// otherwise, only expose an address if it is in the list of exposable addresses
		//
		// TODO: should not expose everything by default...
		return addressMap == nil || addressMap[address]
	}
}

func onConnect(ctx context.Context, _ *remotedialer.Session) error {
	return nil
}
