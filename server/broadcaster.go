package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func NewBroadcaster() Broadcaster {
	return &TCPBroadcaster{
		parentConnection: nil,
		childConnections: make(map[string]net.Conn),
	}
}

type Broadcaster interface {
	Broadcast(cmd string)
	AddChildConnection(conn net.Conn)
	AddParentConnection(conn net.Conn)
	RemoveChildConnection(address string)
	RemoveParentConnection()
}

type TCPBroadcaster struct {
	sync.RWMutex

	parentConnection net.Conn
	childConnections map[string]net.Conn
}

func (b *TCPBroadcaster) Broadcast(cmd string) {
	for address, childConn := range b.childConnections {
		if _, err := fmt.Fprint(childConn, cmd); err != nil {
			b.parentConnection = nil
			log.Println(fmt.Sprintf("error broadcasting message to child node. node is unavailable: %s", address))
		}
	}

	if b.parentConnection == nil {
		return
	}

	if _, err := fmt.Fprint(b.parentConnection, cmd); err != nil {
		b.parentConnection = nil
		log.Println(fmt.Sprintf("error broadcasting message to parent node. node is unavailable: %s", b.parentConnection.RemoteAddr().String()))
		return
	}
}

func (b *TCPBroadcaster) AddChildConnection(conn net.Conn) {
	b.childConnections[conn.RemoteAddr().String()] = conn
	log.Println(fmt.Sprintf("child node has been connected. address %s", conn.RemoteAddr().String()))
}

func (b *TCPBroadcaster) AddParentConnection(conn net.Conn) {
	b.parentConnection = conn
	log.Println(fmt.Sprintf("parent node has been connected. address %s", conn.RemoteAddr().String()))
}

func (b *TCPBroadcaster) RemoveChildConnection(address string) {
	delete(b.childConnections, address)
	log.Println(fmt.Sprintf("child node has been disconnected. address %s", address))
}

func (b *TCPBroadcaster) RemoveParentConnection() {
	b.parentConnection = nil
	log.Println(fmt.Sprintf("parent node has been disconnected"))
}

func (b *TCPBroadcaster) HasChildNodes() bool {

}

func (b *TCPBroadcaster) HasParentNodes() {

}
