package gateway

import "github.com/aiyengar2/portexporter/pkg/config"

// Config represents the configuration of a Gateway
type Config struct {
	config.TLSClient
	Expose []string `yaml:"expose,omitempty"`
}
