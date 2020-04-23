package dic

import (
	"errors"

	"github.com/sarulabs/di"
	"github.com/sarulabs/dingo"

	log "log"

	service "github.com/antelman107/dependency_injection_examples/sarulabs/dingo/service"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "config",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("config")
				if err != nil {
					var eo *service.Config
					return eo, err
				}
				b, ok := d.Build.(func() (*service.Config, error))
				if !ok {
					var eo *service.Config
					return eo, errors.New("could not cast build function to func() (*service.Config, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "logger",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logger")
				if err != nil {
					var eo *log.Logger
					return eo, err
				}
				pi0, err := ctn.SafeGet("config")
				if err != nil {
					var eo *log.Logger
					return eo, err
				}
				p0, ok := pi0.(*service.Config)
				if !ok {
					var eo *log.Logger
					return eo, errors.New("could not cast parameter 0 to *service.Config")
				}
				b, ok := d.Build.(func(*service.Config) (*log.Logger, error))
				if !ok {
					var eo *log.Logger
					return eo, errors.New("could not cast build function to func(*service.Config) (*log.Logger, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
	}
}
