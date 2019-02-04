package server

import (
	"github.com/caarlos0/env"
)

func NewConfiguration() Configuration {
	config := Configuration{}

	if err := env.Parse(&config); err != nil {
		panic("error parsing environment configuration: " + err.Error())
	}

	return config
}

type Configuration struct {
	ServiceHost string `env:"SERVICE_HOST,required"`
	ServicePort string `env:"SERVICE_PORT,required"`
}
