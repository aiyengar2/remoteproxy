package redirect

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/aiyengar2/portexporter/pkg/config"
	"github.com/fsnotify/fsnotify"
)

type Redirect struct {
	config.HTTP
	config.TLSClient
	Address string `yaml:"address,omitempty"`
}

func (r Redirect) ToHandler() http.Handler {
	return &httputil.ReverseProxy{
		Director: r.HTTP.Director,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       r.TLSConfig(r.Address),
		},
	}
}

func (r Redirect) RestartWatcher() (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	for _, path := range []string{r.TokenFile, r.CaCertFile} {
		if path == "" {
			continue
		}
		if err := w.Add(path); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (r Redirect) String() string {
	return fmt.Sprintf("[address=%s,http=%s,tls=%s]", r.Address, r.HTTP, r.TLSClient)
}
