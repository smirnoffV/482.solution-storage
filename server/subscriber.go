package server

import (
	"482.solutions-node-storage/storage"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func NewSubscriber(handler Handler, broadcaster Broadcaster, configuration Configuration, storage *storage.Storage) Subscriber {
	return &ParentNodeSubscriber{
		Handler:       handler,
		Broadcaster:   broadcaster,
		Configuration: configuration,
		Storage:       storage,
		isRecovered:   false,
	}
}

type Subscriber interface {
	Subscribe()
}

type ParentNodeSubscriber struct {
	Handler       Handler
	Broadcaster   Broadcaster
	Configuration Configuration
	Storage       *storage.Storage

	isRecovered bool
}

func (s *ParentNodeSubscriber) Subscribe() {
	if !s.Configuration.IsParentNodeAddressSet() {
		return
	}

	conn, err := net.Dial("tcp", s.Configuration.GetParentNodeAddress())
	if err != nil {
		log.Println("error subscribing to parent node")
		os.Exit(0)
	}

	s.Broadcaster.AddParentConnection(conn)

	if !s.isRecovered {
		cmd := fmt.Sprintf("%s\n", IAmChild)
		if _, err := fmt.Fprint(conn, cmd); err != nil {
			s.Broadcaster.RemoveParentConnection()
			log.Println("error recovering data from parent node")
			os.Exit(0)
		}
	}

	for {
		cmd, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			s.Broadcaster.RemoveParentConnection()
			fmt.Println("error reading response")
			break
		}

		s.Handler.ProcessCommand(cmd, conn)
	}

}
