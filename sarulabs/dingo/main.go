package main

import (
	"log"
	"os"

	"github.com/sarulabs/dingo"

	"github.com/antelman107/dependency_injection_examples/sarulabs/dingo/generated/dic"
	"github.com/antelman107/dependency_injection_examples/sarulabs/dingo/service"
)

func main() {
	if len(os.Args) == 2 {
		err := dingo.GenerateContainer((*service.Provider)(nil), os.Args[1])
		if err != nil {
			panic(err)
		}

		return
	}

	c, err := dic.NewContainer()
	if err != nil {
		panic(err)
	}

	var logger *log.Logger
	err = c.Fill("logger", &logger)
	if err != nil {
		panic(err)
	}

	logger.Println("ok")
}
