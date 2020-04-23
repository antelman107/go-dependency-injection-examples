package service

import (
	"encoding/json"
	"github.com/sarulabs/dingo"
	"log"
	"os"
)

var ServicesADefs = []dingo.Def{
	{
		Name: "config",
		Build: func() (cfg *Config, err error) {
			err = json.Unmarshal([]byte(`{"prefix": "[foo] "}`), &cfg)
			return
		},
	},
	{
		Name: "logger",
		Build: func(cfg *Config) (l *log.Logger, err error) {
			return log.New(os.Stdout, cfg.Prefix, 0), nil
		},
	},
}

type Config struct {
	Prefix string
}

// Redefine your own service by overriding the Load method of the dingo.BaseProvider.
type Provider struct {
	dingo.BaseProvider
}

func (p *Provider) Load() error {
	return p.AddDefSlice(ServicesADefs)
}
