package impl

import (
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"github.com/ryan-berger/chatty/manager/connection"
	"github.com/ryan-berger/chatty/repositories/models"
	"time"
)

type RequestType string

const (
	SendMessage        RequestType = "sendMessage"
	CreateConversation RequestType = "CreateConversation"
)

var stringToType = map[RequestType]connection.RequestType{
	SendMessage:        connection.SendMessage,
	CreateConversation: connection.CreateConversation,
}

type wsConn struct {
	conn       *websocket.Conn
	conversant models.Conversant
	leave      chan struct{}
	requests   chan connection.Request
	responses  chan connection.Response
}

type wsRequest struct {
	RequestType RequestType
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

func newWsConn(conn *websocket.Conn) wsConn {
	wsConn := wsConn{
		conversant: models.Conversant{Id: uuid.New()},
		conn:       conn,
		leave:      make(chan struct{}, 1),
		requests:   make(chan connection.Request),
		responses:  make(chan connection.Response),
	}
	return wsConn
}

func (conn wsConn) pumpIn() {
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

func (conn wsConn) pumpOut() {
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

func (conn wsConn) send(response connection.Response) {
	conn.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := conn.conn.WriteJSON(&response)
	if err != nil {
		conn.leave <- struct{}{}
	}
}

func (conn wsConn) receive() {
	var req wsRequest
	conn.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	err := conn.conn.ReadJSON(&req)
	if err != nil {
		conn.leave <- struct{}{}
	}
	conn.requests <- wsRequestToRequest(req)
}

func (conn wsConn) Authorize() error {
	err := conn.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		conn.conn.Close()
		return err
	}
	var test map[string]string
	err = conn.conn.ReadJSON(&test)

	if err != nil {
		conn.conn.Close()
		return err
	}

	go conn.pumpIn()
	go conn.pumpOut()
	return nil
}

func (conn wsConn) GetConversant() models.Conversant {
	return conn.conversant
}

func (conn wsConn) Requests() chan connection.Request {
	return conn.requests
}

func (conn wsConn) Response() chan connection.Response {
	return conn.responses
}

func (conn wsConn) Leave() chan struct{} {
	return conn.leave
}
