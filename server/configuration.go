package server

import (
	"github.com/caarlos0/env"
	"net"
)

func NewConfiguration() Configuration {
	config := Configuration{}

	if err := env.Parse(&config); err != nil {
		panic("error parsing environment configuration: " + err.Error())
	}

	return config
}

type Configuration struct {
	ServiceHost           string `env:"SERVICE_HOST,required"`
	ServicePort           string `env:"SERVICE_PORT,required"`
	ParentNodeServiceHost string `env:"PARENT_NODE_SERVICE_HOST,required"`
	ParentNodeServicePort string `env:"PARENT_NODE_SERVICE_PORT,required"`
}

func (c Configuration) IsParentNodeAddressSet() bool {
	return c.GetParentNodeAddress() != ":"
}

func (c Configuration) GetParentNodeAddress() string {
	address := net.JoinHostPort(c.ParentNodeServiceHost, c.ParentNodeServicePort)
	return address
}
