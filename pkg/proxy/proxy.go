package proxy

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/aiyengar2/portexporter/pkg/config"
	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
)

type proxyServer struct {
	http.Server

	useTLS bool
}

func NewServer(listenAddr string, config config.TLSServer) *proxyServer {
	s := &proxyServer{}

	if config.CertFile != "" && config.KeyFile != "" {
		s.useTLS = true
	}

	authorizer := func(req *http.Request) (string, bool, error) {
		id := req.Header.Get("X-Proxy-Tunnel-ID")
		return id, id != "", nil
	}
	s.Server = http.Server{
		Addr:         listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		// disable HTTP/2 support
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		Handler: &proxyHandler{
			rdServer: remotedialer.New(authorizer, remotedialer.DefaultErrorWriter),
		},
		TLSConfig: config.TLSConfig(listenAddr),
	}

	return s
}

func (s *proxyServer) Start(ctx context.Context) error {
	go func() {
		if !s.useTLS {
			logrus.Infof("Listening for HTTP connections on %s", s.Addr)
			if err := s.ListenAndServe(); err != nil {
				logrus.Error(err)
			}
		} else {
			logrus.Infof("Listening for TLS connections on %s", s.Addr)
			if err := s.ListenAndServeTLS("", ""); err != nil {
				logrus.Error(err)
			}
		}
	}()
	<-ctx.Done()
	logrus.Infof("Shutting down...")
	return s.Shutdown(ctx)
}
