package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rancher/remotedialer"
	"github.com/sirupsen/logrus"
)

type proxyHandler struct {
	rdServer *remotedialer.Server
}

func (h *proxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logrus.Debugf("Received request from host [%s] to url [%s] for method %s", req.RemoteAddr, req.URL, req.Method)
	if req.URL.Host == "" {
		if req.URL.Path == "/connect" {
			h.rdServer.ServeHTTP(rw, req)
			return
		}
		http.Error(rw, "proxy only supports '/connect'", http.StatusNotFound)
		return
	}
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
	if req.Method == http.MethodConnect {
		h.handleHTTPS(rw, req)
		return
	}
	h.handleHTTP(rw, req)
}

func (h *proxyHandler) handleHTTPS(rw http.ResponseWriter, req *http.Request) {
	tunnelConn, err := h.getDialer(req)(context.TODO(), "tcp", req.Host)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
	rw.WriteHeader(http.StatusOK)

	// hijack incoming HTTPS connection
	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "connection does not support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw, fmt.Sprintf("cannot hijack connection: %s", err), http.StatusInternalServerError)
	}

	pipe := func(dst io.WriteCloser, src io.ReadCloser) {
		defer dst.Close()
		defer src.Close()
		io.Copy(dst, src)
	}

	go pipe(tunnelConn, conn)
	go pipe(conn, tunnelConn)
}

func (h *proxyHandler) handleHTTP(rw http.ResponseWriter, req *http.Request) {
	// send packets over the wire and wait for a response
	transport := http.Transport{
		DialContext: h.getDialer(req),
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// pipe response
	rwHeader := rw.Header()
	for k, vSlice := range resp.Header {
		for _, v := range vSlice {
			rwHeader.Add(k, v)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	io.Copy(rw, resp.Body)
}

func (h *proxyHandler) getDialer(req *http.Request) remotedialer.Dialer {
	return h.rdServer.Dialer(strings.Split(req.Host, ":")[0])
}
