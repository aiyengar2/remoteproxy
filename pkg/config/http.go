package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type HTTP struct {
	TokenFile string `yaml:"tokenFile,omitempty"`

	token    []byte
	readOnce sync.Once
}

func (h HTTP) String() string {
	return fmt.Sprintf("[tokenFile=%s]", h.TokenFile)
}

func (h *HTTP) Director(req *http.Request) {
	if h.TokenFile == "" {
		req.URL.Scheme = "http"
		return
	}
	h.readOnce.Do(h.setToken)
	req.URL.Scheme = "https"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.token))
}

func (h *HTTP) setToken() {
	var err error
	h.token, err = ioutil.ReadFile(h.TokenFile)
	if err != nil {
		logrus.Warnf("could not read token from path %s", h.TokenFile)
	}
	if len(h.token) == 0 {
		logrus.Warnf("no token found at path %s", h.TokenFile)
	}
}
