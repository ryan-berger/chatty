package websocket

import (
	"errors"
	"time"

	"github.com/ryan-berger/chatty/repositories"

	"github.com/ryan-berger/chatty/connection"

	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
)

type requestType string

const (
	sendMessage        requestType = "sendMessage"
	createConversation requestType = "createConversation"
)

var stringToType = map[requestType]connection.RequestType{
	sendMessage:        connection.SendMessage,
	createConversation: connection.CreateConversation,
}

type Auth func(map[string]string) (repositories.Conversant, error)

// Conn is a websocket implementation of the Conn interface
type Conn struct {
	conn       *websocket.Conn
	conversant repositories.Conversant
	leave      chan struct{}
	requests   chan connection.Request
	responses  chan connection.Response
	auth       Auth
}

type wsRequest struct {
	RequestType requestType
	Data        interface{}
}

func wsRequestToRequest(request wsRequest) connection.Request {
	req := connection.Request{}

	if _, ok := stringToType[request.RequestType]; !ok {
		return req
	}

	req.Type = stringToType[request.RequestType]
	return req
}

// NewWebsocketConn is a factory for a websocket connection
func NewWebsocketConn(conn *websocket.Conn, auth Auth) Conn {
	wsConn := Conn{
		conversant: repositories.Conversant{ID: uuid.New()},
		conn:       conn,
		leave:      make(chan struct{}, 1),
		requests:   make(chan connection.Request),
		responses:  make(chan connection.Response),
		auth:       auth,
	}
	return wsConn
}

func (conn Conn) pumpIn() {
	for {
		select {
		case <-conn.leave:
			conn.conn.Close()
			return
		default:
			conn.receive()
		}
	}
}

func (conn Conn) pumpOut() {
	for {
		select {
		case <-conn.leave:
			conn.conn.Close()
			return
		case response := <-conn.responses:
			conn.send(response)
		}
	}
}

func (conn Conn) send(response connection.Response) {
	conn.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := conn.conn.WriteJSON(&response)
	if err != nil {
		conn.leave <- struct{}{}
	}
}

func (conn Conn) receive() {
	var req wsRequest
	conn.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	err := conn.conn.ReadJSON(&req)
	if err != nil {
		conn.leave <- struct{}{}
	}
	conn.requests <- wsRequestToRequest(req)
}

// Authorize satisfies the Conn interface
func (conn Conn) Authorize() error {
	err := conn.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		conn.conn.Close()
		return err
	}
	var creds map[string]string
	err = conn.conn.ReadJSON(&creds)

	if err != nil {
		conn.conn.Close()
		return err
	}
	conversant, err := conn.auth(creds)

	if err != nil {
		conn.conn.Close()
		return errors.New("not authorized")
	}

	conn.conversant = conversant

	go conn.pumpIn()
	go conn.pumpOut()
	return nil
}

// GetConversant satisfies the Conn interface
func (conn Conn) GetConversant() repositories.Conversant {
	return conn.conversant
}

// Requests satisfies the Conn interface
func (conn Conn) Requests() chan connection.Request {
	return conn.requests
}

// Response satisfies the Conn interface
func (conn Conn) Response() chan connection.Response {
	return conn.responses
}

// Leave satisfies the Conn interface
func (conn Conn) Leave() chan struct{} {
	return conn.leave
}
