// Code generated by dingo; DO NOT EDIT
package main

import (
	log "log"
	"os"
)

type Container struct {
	Config *Config
	Logger *log.Logger
}

var DefaultContainer = NewContainer()

func NewContainer() *Container {
	return &Container{}
}
func (container *Container) GetConfig() *Config {
	if container.Config == nil {
		service, err := NewConfig()
		if err != nil {
			return nil
		}
		container.Config = service
	}
	return container.Config
}
func (container *Container) GetLogger() *log.Logger {
	if container.Logger == nil {
		service := log.New(os.Stdout, container.GetConfig().Prefix, 0)
		container.Logger = service
	}
	return container.Logger
}
