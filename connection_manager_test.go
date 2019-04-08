package chatty

import (
	"sync"
	"testing"
	"time"

	"github.com/pborman/uuid"

	"github.com/ryan-berger/chatty/connection"
	"github.com/ryan-berger/chatty/operators"

	"github.com/ryan-berger/chatty/repositories"
)

type testData struct {
	*repositories.MockConversationRepo
	*repositories.MockConversantRepo
	*repositories.MockMessageRepo
	*operators.MockAuther
	*operators.MockNotifier
}

func newTestData() *testData {
	return &testData{
		&repositories.MockConversationRepo{},
		&repositories.MockConversantRepo{},
		&repositories.MockMessageRepo{},
		&operators.MockAuther{},
		&operators.MockNotifier{},
	}
}

func makeConn(id string) *connection.MockConn {
	mockConn := &connection.MockConn{}
	mockConn.Conversant = func() repositories.Conversant {
		return repositories.Conversant{
			ID: id,
		}
	}
	mockConn.Auth = func() error {
		return nil
	}
	mockConn.Request = func() chan connection.Request {
		return nil
	}
	mockConn.Leaver = func() chan struct{} {
		return nil
	}
	return mockConn
}

func makeMockManager() *ConnectionManager {
	return &ConnectionManager{
		connections:    make(map[string][]connection.Conn),
		connectionMu:   &sync.RWMutex{},
		messageChan:    make(chan messageRequest, 10),
		shutdownChan:   make(chan struct{}, 1),
		chatInteractor: &chatInteractor{},
	}
}

func TestConnectionManager_JoinAuthorize(t *testing.T) {
	manager := makeMockManager()
	manager.startup()

	conversantRepo := &repositories.MockConversantRepo{
		Upsert: func(conversant repositories.Conversant) (*repositories.Conversant, error) {
			return &conversant, nil
		},
	}

	manager.chatInteractor.conversantRepo = conversantRepo
	conn := makeConn(uuid.New())

	manager.Join(conn)

	if len(manager.connections) != 1 {
		t.Fatalf("Connections didn't join")
	}

}

func TestConnectionManager_NotifyInMemory(t *testing.T) {
	manager := makeMockManager()
	manager.startup()

	senderID := uuid.New()
	receiverID := uuid.New()

	m := &repositories.MockMessageRepo{
		Create: func(message repositories.Message) (message2 *repositories.Message, e error) {
			return &message, nil
		},
	}

	c := &repositories.MockConversationRepo{
		GetConvo: func(conversationId string) (conversants []repositories.Conversant, e error) {
			return []repositories.Conversant{
				{ID: receiverID},
				{ID: senderID},
			}, nil
		},
	}

	manager.chatInteractor = newChatInteractor(m, c, nil)

	connA := makeConn(senderID)
	connB := makeConn(receiverID)

	resp := make(chan connection.Response, 2)

	connA.Resp = func() chan connection.Response {
		return resp
	}

	connB.Resp = func() chan connection.Response {
		return resp
	}

	manager.addConn(connA)
	manager.addConn(connB)

	manager.sendMessage(connA, connection.SendMessageRequest{ConversationID: uuid.New(), Message: "test"})

	select {
	case response := <-resp:
		if response.Type == connection.Error {
			t.Fatalf("received response error")
		}
		return
	case <-time.After(10 * time.Second):
		t.Error("Didn't receive message")
	}
}

func TestConnectionManager_Notifier(t *testing.T) {
	manager := makeMockManager()
	manager.startup()

	m := &repositories.MockMessageRepo{
		Create: func(message repositories.Message) (message2 *repositories.Message, e error) {
			return &message, nil
		},
	}

	c := &repositories.MockConversationRepo{
		GetConvo: func(conversationId string) (conversants []repositories.Conversant, e error) {
			return []repositories.Conversant{
				{ID: "b"},
			}, nil
		},
	}
	manager.chatInteractor = newChatInteractor(m, c, nil)
	conn := makeConn("a")

	notified := make(chan struct{}, 1)
	requests := make(chan connection.Request)

	conn.Request = func() chan connection.Request {
		return requests
	}

	conn.Resp = func() chan connection.Response {
		return make(chan connection.Response)
	}

	manager.notifier = &operators.MockNotifier{
		SendNotification: func(id string, message repositories.Message) error {
			notified <- struct{}{}
			return nil
		},
	}

	manager.addConn(conn)

	requests <- connection.Request{Type: connection.SendMessage, Data: connection.SendMessageRequest{ConversationID: uuid.New(), Message: "Test"}}

	select {
	case <-notified:
		return
	case <-time.After(10 * time.Millisecond):
		t.Error("didn't receive notification")
	}
}

func TestManager_Leave(t *testing.T) {
	td := newTestData()
	manager := NewManager(td, td, td, td, td)

	manager.chatInteractor.conversantRepo = &repositories.MockConversantRepo{
		Upsert: func(conversant repositories.Conversant) (*repositories.Conversant, error) {
			return &conversant, nil
		},
	}

	conn1 := makeConn("a")

	closer := make(chan struct{})

	conn1.Leaver = func() chan struct{} {
		return closer
	}

	conn1.Resp = func() chan connection.Response {
		return nil
	}

	conn1.Request = func() chan connection.Request {
		return nil
	}

	conn1.Auth = func() error {
		return nil
	}

	manager.Join(conn1)

	if len(manager.connections) == 0 {
		t.Error("conn didn't join")
	}

	closer <- struct{}{}
	manager.connectionMu.Lock()
	if len(manager.connections) != 0 {
		t.Error("conn didn't leave")
	}
	manager.connectionMu.Unlock()
}
