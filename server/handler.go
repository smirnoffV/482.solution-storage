package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	WrongParamsFormatErr = errors.New("wrong params format")
)

const (
	GETPrefix     = "GET"
	SETPrefix     = "SET"
	GETALLPrefix  = "GETALL"
	IAmChild      = "I_AM_CHILD"
	RECOVERPrefix = "RECOVER"
)

func NewHandler(api Api, listener net.Listener, broadcaster Broadcaster, commandsChan chan string) Handler {
	return &TCPHandler{
		Api:          api,
		TcpListener:  listener,
		Broadcaster:  broadcaster,
		CommandsChan: commandsChan,
	}
}

type Handler interface {
	Handle() error
	HandleConn(net.Conn)
	ProcessCommand(cmd string, conn net.Conn)
}

type TCPHandler struct {
	Api         Api
	TcpListener net.Listener
	Broadcaster Broadcaster

	CommandsChan chan string
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

	for {

		cmd, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Print("client has been disconnected")
			break
		}

		h.ProcessCommand(cmd, conn)
	}

	return
}

func (h *TCPHandler) ProcessCommand(cmd string, conn net.Conn) {
	request := NewRawRequest(cmd)

	var err error
	response := new(bytes.Buffer)

	switch {
	case request.Method == SETPrefix && !request.IsEmptyBody():
		response, err = h.Api.Set(request.Body)
		h.Broadcaster.Broadcast(request.BuildCmd())
	case request.Method == GETPrefix && !request.IsEmptyBody():
		response, err = h.Api.Get(request.Body)
	case request.Method == IAmChild:
		h.Broadcaster.AddChildConnection(conn)
		response, err = h.Api.BuildRecoverResponse()
	case request.Method == RECOVERPrefix && !request.IsEmptyBody():
		h.Api.Recover(request.Body)
	default:
		err = WrongParamsFormatErr
	}

	if err != nil {
		sendErrorResponse(err, conn)
		return
	}

	response.WriteString("\n")

	if _, err := response.WriteTo(conn); err != nil {
		log.Println("error publishing response")
	}
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
	message = strings.Replace(message, "\r", "", -1)
	message = strings.Replace(message, "\n", "", -1)

	parted := strings.Split(message, "||")

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

func (r RawRequest) BuildCmd() string {
	return fmt.Sprintf("%s||%s\n", r.Method, r.Body)
}

func (r RawRequest) IsEmptyBody() bool {
	return r.Body == ""
}
