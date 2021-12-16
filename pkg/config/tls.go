package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type TLSServer struct {
	CertFile   string `yaml:"certFile,omitempty"`
	KeyFile    string `yaml:"keyFile,omitempty"`
	CaCertFile string `yaml:"caCertFile,omitempty"`
}

func (s TLSServer) TLSConfig(address string) *tls.Config {
	tlsConfig := &tls.Config{
		ServerName: address,
	}
	err := s.loadConfig(tlsConfig)
	if err != nil {
		logrus.Error(err)
	}
	return tlsConfig
}

func (s TLSServer) loadConfig(tlsConfig *tls.Config) error {
	if s.CertFile == "" || s.KeyFile == "" {
		return nil
	}
	cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
	if err != nil {
		return fmt.Errorf("unable to load X.509 certificate from cert file %s and key file %s: %s", s.CertFile, s.KeyFile, err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	if s.CaCertFile == "" {
		return nil
	}
	caCert, err := ioutil.ReadFile(s.CaCertFile)
	if err != nil {
		return fmt.Errorf("unable to read cacert file %s: %s", s.CaCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return fmt.Errorf("failed to use cacert file %s as ca certificate", s.CaCertFile)
	}
	tlsConfig.ClientCAs = caCertPool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	return nil
}

func (s TLSServer) String() string {
	return fmt.Sprintf("[certFile=%s,keyFile=%s]", s.CertFile, s.KeyFile)
}

type TLSClient struct {
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify,omitempty"`
	CaCertFile         string `yaml:"caCertFile,omitempty"`
}

func (c TLSClient) TLSConfig(address string) *tls.Config {
	tlsConfig := &tls.Config{
		ServerName: address,
	}
	err := c.loadConfig(tlsConfig)
	if err != nil {
		logrus.Error(err)
	}
	return tlsConfig
}

func (c TLSClient) loadConfig(tlsConfig *tls.Config) error {
	tlsConfig.InsecureSkipVerify = c.InsecureSkipVerify
	if c.CaCertFile == "" {
		return nil
	}
	caCert, err := ioutil.ReadFile(c.CaCertFile)
	if err != nil {
		return fmt.Errorf("unable to read cacert file %s: %s", c.CaCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return fmt.Errorf("failed to use cacert file %s as ca certificate", c.CaCertFile)
	}
	tlsConfig.RootCAs = caCertPool

	return nil
}

func (c TLSClient) String() string {
	return fmt.Sprintf("[insecureSkipVerify=%t,caCertFile=%s]", c.InsecureSkipVerify, c.CaCertFile)
}
