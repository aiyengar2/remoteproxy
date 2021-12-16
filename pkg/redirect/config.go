package redirect

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	DefaultRedirectConfigFile = "redirect.yaml"
)

type Config struct {
	Redirect []Redirect `yaml:"redirect,omitempty"`
}

func Load(configFile string) (Config, error) {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	var opts Config
	return opts, yaml.Unmarshal(configBytes, &opts)
}
