package redirect

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type redirectServer struct {
	*http.Server
}

func NewServer(listenAddr string, config Config) *redirectServer {
	s := &redirectServer{}
	router := Router()
	for _, redirect := range config.Redirect {
		router.RegisterHandler(redirect.Address, redirect)
	}
	s.Server = &http.Server{
		Addr:         listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		// disable HTTP/2 support
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		Handler:      router,
	}
	return s
}

func (s *redirectServer) Start(ctx context.Context) error {
	logrus.Infof("Listening on %s", s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
	<-ctx.Done()
	logrus.Infof("Shutting down...")
	return s.Shutdown(ctx)
}
