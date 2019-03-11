package implementations

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ryan-berger/chatty/repositories"

	"github.com/ryan-berger/chatty/connection"
)

type testConn struct {
	readChan  chan []byte
	writeChan chan []byte
	isClosed  bool
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

func (conn *testConn) WriteJSON(in interface{}) error {
	b, _ := json.Marshal(in)
	conn.writeChan <- b
	return nil
}

func (*testConn) ReadJSON(interface{}) error {
	return nil
}

func (conn *testConn) Close() error {
	conn.isClosed = true
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

var responses = []struct {
	respType connection.ResponseType
	resp     []byte
}{
	{
		respType: connection.NewMessage,
		resp:     []byte(`{"type":"newMessage","data":null}`),
	},
	{
		respType: connection.NewConversation,
		resp:     []byte(`{"type":"newConversation","data":null}`),
	},
	{
		respType: connection.Error,
		resp:     []byte(`{"type":"error","data":null}`),
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
		req := <-conn.Requests()

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

func TestConn_Responses(t *testing.T) {
	responseChan := make(chan connection.Response)
	writeChan := make(chan []byte)

	conn := Conn{
		conn: &testConn{
			writeChan: writeChan,
		},
		leave:     make(chan struct{}, 1),
		responses: responseChan,
	}

	go conn.pumpOut()

	for _, response := range responses {
		conn.Response() <- connection.Response{Type: response.respType}

		res := <-writeChan

		if string(res) != string(response.resp) {
			t.Fatalf("expected %s, received %s", string(response.resp), string(res))
		}
	}

	conn.Leave() <- struct{}{}
}

func TestConn_Authorize(t *testing.T) {
	readChan := make(chan []byte, 1)

	testConn := &testConn{
		readChan: readChan,
	}

	conn := NewWebsocketConn(testConn, func(strings map[string]string) (conversant repositories.Conversant, e error) {
		return repositories.Conversant{}, errors.New("test")
	})
	readChan <- []byte(`{"test": "test"}`)
	conn.Authorize()

	if testConn.isClosed != true {
		t.Fatalf("connection should have closed")
	}
}

func TestConn_GetConversant(t *testing.T) {
	readChan := make(chan []byte, 1)

	testConn := &testConn{
		readChan: readChan,
	}

	conn := NewWebsocketConn(testConn, func(strings map[string]string) (conversant repositories.Conversant, e error) {
		return repositories.Conversant{ID: "testID"}, nil
	})

	readChan <- []byte(`{"test": "test"}`)
	conn.Authorize()

	if conn.GetConversant().ID != "testID" {
		t.Fatalf("expected that %s would be testID", conn.GetConversant().ID)
	}
}
