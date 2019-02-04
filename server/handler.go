package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/dig"
	"log"
	"net"
	"strings"
)

const (
	GETPrefix    = "GET"
	SETPrefix    = "SET"
	GETALLPrefix = "GETALL"
)

func NewHandler(deps TCPHandlerDependencies) Handler {
	return &TCPHandler{
		deps,
	}
}

type Handler interface {
	Handle() error
	HandleConn(net.Conn)
}

type TCPHandlerDependencies struct {
	dig.In

	Api         Api
	TcpListener net.Listener `name:"request-server"`
}

type TCPHandler struct {
	TCPHandlerDependencies
}

func (h *TCPHandler) Handle() error {
	for {
		conn, err := h.TcpListener.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}

		go h.HandleConn(conn)
	}

	return nil
}

func (h *TCPHandler) HandleConn(conn net.Conn) {
	defer conn.Close()

	log.Println("Goroutine has been started")

	for {

		message, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			log.Print("client has been disconnected")
			break
		}

		request := NewRawRequest(message)

		var response *bytes.Buffer

		switch {
		case request.Method == SETPrefix && !request.IsEmptyBody():
			response, err = h.Api.Set(request.Body)
		case request.Method == GETPrefix && !request.IsEmptyBody():
			response, err = h.Api.Get(request.Body)
		case request.Method == GETALLPrefix:
			response, err = h.Api.GetAll()
		default:
			err = errors.New("wrong params format")
		}

		if err != nil {
			sendErrorResponse(err, conn)
			continue
		}

		response.WriteString("\n")

		if _, err := response.WriteTo(conn); err != nil {
			log.Println("error publishing response")
		}

	}

	return
}

func sendErrorResponse(err error, conn net.Conn) {
	responseBytes, err := json.Marshal(ErrorResponse{Error: err.Error()})
	if err != nil {
		log.Println("error marshalling response")
		return
	}

	buffer := bytes.NewBuffer(responseBytes)
	buffer.WriteString("\n")
	if _, err := buffer.WriteTo(conn); err != nil {
		log.Println("error publishing response")
		return
	}
}

func NewRawRequest(message string) RawRequest {
	parted := strings.Split(strings.TrimSuffix(message, "\r\n"), "||")

	request := RawRequest{
		Method: parted[0],
	}

	if len(parted) == 2 {
		request.Body = parted[1]
	}

	return request
}

type RawRequest struct {
	Method string
	Body   string
}

func (r RawRequest) IsEmptyBody() bool {
	return r.Body == ""
}
