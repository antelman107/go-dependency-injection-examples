package main

import (
	"encoding/json"
	"go.uber.org/dig"
	"log"
	"os"
)

type Config struct {
	Prefix string
}

func main() {
	c := dig.New()
	err := c.Provide(func() (*Config, error) {
		// In a real program, the configuration will probably be read from a
		// file.
		var cfg Config
		err := json.Unmarshal([]byte(`{"prefix": "[foo] "}`), &cfg)
		return &cfg, err
	})
	if err != nil {
		panic(err)
	}

	// Provide a way to build the logger based on the configuration.
	err = c.Provide(func(cfg *Config) *log.Logger {
		return log.New(os.Stdout, cfg.Prefix, 0)
	})
	if err != nil {
		panic(err)
	}

	// The second call with same paramers will cause a panic.
	// Use name parameter to fix it
	err = c.Provide(func(cfg *Config) *log.Logger {
		return log.New(os.Stdout, cfg.Prefix, 0)
	}, dig.Name("logger2"))
	if err != nil {
		panic(err)
	}

	// Invoke a function that requires the logger, which in turn builds the
	// Config first.
	err = c.Invoke(func(l *log.Logger) {
		l.Print("You've been invoked")
	})
	if err != nil {
		panic(err)
	}
}
