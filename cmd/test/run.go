package test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aiyengar2/portexporter/pkg/config"
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
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
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

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cfg := config.TLSServer{
		CertFile:   certFile,
		KeyFile:    keyFile,
		CaCertFile: caCertFile,
	}

	server := http.Server{
		Addr:      listen,
		Handler:   &testHandler{},
		TLSConfig: cfg.TLSConfig(listen),
	}
	go func() {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			logrus.Infof("Listening for HTTP requests on http://%s", listen)
			if err := server.ListenAndServe(); err != nil {
				logrus.Error(err)
			}
		} else {
			logrus.Infof("Listening for TLS connections on %s", listen)
			if err := server.ListenAndServeTLS("", ""); err != nil {
				logrus.Error(err)
			}
		}
	}()
	<-ctx.Done()
	logrus.Infof("Shutting down...")
	return server.Shutdown(ctx)
}

type testHandler struct{}

func (h *testHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logrus.Debugf("Received request from host [%s] to url [%s]", req.RemoteAddr, req.URL)
	var message string
	if req.URL.Path != "/" {
		message = strings.ReplaceAll(req.URL.Path, "/", " ")
	}
	fmt.Fprintf(rw, "hello%s", message)
}
