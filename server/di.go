package server

import (
	"482.solutions-node-storage/storage"
	"fmt"
	"go.uber.org/dig"
	"log"
	"net"
)

func NewContainer() *dig.Container {
	c := dig.New()

	if err := c.Provide(func() *storage.Storage {
		s := new(storage.Storage)
		s.Data = make(map[string]string)
		return s
	}); err != nil {
		panic(err)
	}

	if err := c.Provide(NewConfiguration); err != nil {
		panic(err)
	}

	if err := c.Provide(NewApi); err != nil {
		panic(err)
	}

	if err := c.Provide(storage.NewRepository); err != nil {
		panic(err)
	}

	if err := c.Provide(NewHandler); err != nil {
		panic(err)
	}

	if err := c.Provide(NewTcpRequestServer); err != nil {
		panic(err)
	}

	if err := c.Provide(NewCommandsChanel); err != nil {
		panic(err)
	}

	if err := c.Provide(NewBroadcaster); err != nil {
		panic(err)
	}

	if err := c.Provide(NewSubscriber); err != nil {
		panic(err)
	}

	return c
}

func NewTcpRequestServer(configuration Configuration) (net.Listener, error) {
	log.Println(fmt.Sprintf("Launching TCP listener... IP:%s PORT:%s", configuration.ServiceHost, configuration.ServicePort))
	ln, err := net.Listen("tcp", net.JoinHostPort(configuration.ServiceHost, configuration.ServicePort))

	if err != nil {
		return nil, err
	}

	return ln, nil
}

func NewCommandsChanel() chan string {
	return make(chan string, 500)
}
