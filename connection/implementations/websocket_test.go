package implementations

import (
	"reflect"
	"testing"
	"time"

	"github.com/ryan-berger/chatty/connection"
)

type testConn struct {
	readChan chan []byte
}

func (*testConn) SetReadDeadline(time.Time) error {
	return nil
}

func (*testConn) SetWriteDeadline(time.Time) error {
	return nil
}

func (conn *testConn) ReadMessage() (int, []byte, error) {
	return 0, <-conn.readChan, nil
}

func (*testConn) WriteJSON(interface{}) error {
	return nil
}

func (*testConn) ReadJSON(interface{}) error {
	return nil
}

func (*testConn) Close() error {
	return nil
}

var requests = []struct {
	reqType connection.RequestType
	reqData interface{}
	req     []byte
}{
	{
		req:     []byte(`{"type": "createConversation"}`),
		reqData: connection.CreateConversationRequest{},
		reqType: connection.CreateConversation,
	},
	{
		req:     []byte(`{"type": "sendMessage"}`),
		reqData: connection.SendMessageRequest{},
		reqType: connection.SendMessage,
	},
	{
		req:     []byte(`{"type": "asdfasdfasdf"}`),
		reqData: nil,
		reqType: connection.RequestError,
	},
}

func TestConn_Requests(t *testing.T) {
	readChan := make(chan []byte)
	requestChan := make(chan connection.Request)

	conn := Conn{
		conn: &testConn{
			readChan: readChan,
		},
		leave:    make(chan struct{}, 1),
		requests: requestChan,
	}

	go conn.pumpIn()

	for _, request := range requests {
		readChan <- request.req
		req := <-requestChan

		if req.Type != request.reqType {
			t.Fatalf("expected %d, type was actually: %d", request.reqType, req.Type)
		}

		if request.reqData == nil && req.Data != nil {
			t.Fatalf("reqData should have been null")
		}

		if request.reqData == nil && req.Data == nil {
			continue
		}

		actual := reflect.TypeOf(req.Data).Name()
		expected := reflect.TypeOf(request.reqData).Name()

		if actual != expected {
			t.Fatalf("expected data type %s, received %s", expected, actual)
		}
	}

	conn.Leave() <- struct{}{}
}
