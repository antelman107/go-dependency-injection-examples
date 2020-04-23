package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	type Config struct {
		Prefix string
	}

	var cfg Config
	err := json.Unmarshal([]byte(`{"prefix": "[foo] "}`), &cfg)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, cfg.Prefix, 0)
	logger.Println("ok")

}
