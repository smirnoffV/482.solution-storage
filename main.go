package main

import (
	"482.solutions-node-storage/server"
)

func main() {
	c := server.NewContainer()

	var err error

	err = c.Invoke(func(subscriber server.Subscriber) {
		go subscriber.Subscribe()
	})

	if err != nil {
		panic(err)
	}

	err = c.Invoke(func(handler server.Handler) {
		handler.Handle()
	})

	if err != nil {
		panic(err)
	}

	return
}
