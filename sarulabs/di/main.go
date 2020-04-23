package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sarulabs/di"
)

type Config struct {
	Prefix string
}

func main() {
	builder, err := di.NewBuilder()
	if err != nil {
		panic(err)
	}

	err = builder.Add(di.Def{
		Name: "config",
		Build: func(ctn di.Container) (interface{}, error) {
			var cfg Config
			err := json.Unmarshal([]byte(`{"prefix": "[foo] "}`), &cfg)
			return &cfg, err
		},
	})
	if err != nil {
		panic(err)
	}

	err = builder.Add(di.Def{
		Name: "logger",
		Build: func(ctn di.Container) (interface{}, error) {
			var cfg *Config
			err = ctn.Fill("config", &cfg)
			return log.New(os.Stdout, cfg.Prefix, 0), nil
		},
		Close: func(obj interface{}) error {
			if _, ok := obj.(*log.Logger); ok {
				fmt.Println("logger close")
			}
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	container := builder.Build()

	var l *log.Logger
	err = container.Fill("logger", &l)
	if err != nil {
		panic(err)
	}

	l.Println("ok")

	err = container.DeleteWithSubContainers()
	if err != nil {
		panic(err)
	}

}
