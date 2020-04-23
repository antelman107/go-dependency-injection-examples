package main

import "encoding/json"

type Config struct {
	Prefix string
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := json.Unmarshal([]byte(`{"prefix": "[foo] "}`), &cfg)
	return &cfg, err
}

func main() {
	DefaultContainer.GetLogger().Println("ok")
}
