package implementations

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ryan-berger/chatty/repositories"

	"github.com/ryan-berger/chatty/connection"

	"github.com/gorilla/websocket"
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
	conn       WebsocketConn
	conversant repositories.Conversant
	leave      chan struct{}
	requests   chan connection.Request
	responses  chan connection.Response
	auth       Auth
}

type WebsocketConn interface {
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	ReadMessage() (int, []byte, error)
	WriteJSON(interface{}) error
	ReadJSON(interface{}) error
	Close() error
}

type wsRequest struct {
	RequestType requestType     `json:"type"`
	Data        json.RawMessage `json:"data"`
}

func wsRequestType(reqType requestType) connection.RequestType {
	if val, ok := stringToType[reqType]; ok {
		return val
	}

	return connection.RequestError
}

func wsRequestData(data []byte) connection.Request {
	var request wsRequest
	json.Unmarshal(data, &request)
	req := connection.Request{Type: wsRequestType(request.RequestType)}
	switch req.Type {
	case connection.SendMessage:
		messageRequest := connection.SendMessageRequest{}
		json.Unmarshal(request.Data, &messageRequest)
		req.Data = messageRequest
	case connection.CreateConversation:
		conversationRequest := connection.CreateConversationRequest{}
		json.Unmarshal(request.Data, &conversationRequest)
		req.Data = conversationRequest
	case connection.RequestError:
		req.Data = nil
	}

	return req
}

// NewWebsocketConn is a factory for a websocket connection
func NewWebsocketConn(conn *websocket.Conn, auth Auth) *Conn {
	wsConn := &Conn{
		conn:      conn,
		leave:     make(chan struct{}, 1),
		requests:  make(chan connection.Request),
		responses: make(chan connection.Response),
		auth:      auth,
	}
	return wsConn
}

func (conn *Conn) pumpIn() {
	conn.conn.SetReadDeadline(time.Time{})
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

func (conn *Conn) pumpOut() {
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

func (conn *Conn) send(response connection.Response) {
	conn.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := conn.conn.WriteJSON(&response)
	if err != nil {
		conn.leave <- struct{}{}
	}
}

func (conn *Conn) receive() {
	_, message, err := conn.conn.ReadMessage()
	if err != nil {
		conn.leave <- struct{}{}
	}
	conn.requests <- wsRequestData(message)
}

// Authorize satisfies the Conn interface
func (conn *Conn) Authorize() error {
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
func (conn *Conn) GetConversant() repositories.Conversant {
	return conn.conversant
}

// Requests satisfies the Conn interface
func (conn *Conn) Requests() chan connection.Request {
	return conn.requests
}

// Response satisfies the Conn interface
func (conn *Conn) Response() chan connection.Response {
	return conn.responses
}

// Leave satisfies the Conn interface
func (conn *Conn) Leave() chan struct{} {
	return conn.leave
}
