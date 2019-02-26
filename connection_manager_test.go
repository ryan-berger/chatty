package chatty

import (
	"testing"
	"time"

	"github.com/ryan-berger/chatty/connection"
	"github.com/ryan-berger/chatty/operators"

	"github.com/ryan-berger/chatty/repositories"
)

type testData struct {
	*repositories.MockConversationRepo
	*repositories.MockMessageRepo
	*operators.MockAuther
	*operators.MockNotifier
}

func newTestData() *testData {
	return &testData{
		&repositories.MockConversationRepo{},
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
	mockConn.Request = func() chan connection.Request {
		return nil
	}
	mockConn.Leaver = func() chan struct{} {
		return nil
	}
	return mockConn
}

func inMemoryMockCons(sender chan struct{}) (conn1, conn2 *connection.MockConn) {
	mockConn1 := makeConn("a")
	mockConn2 := makeConn("b")

	mockConn2.Request = func() chan connection.Request {
		<-sender
		req := make(chan connection.Request, 1)
		req <- connection.Request{Type: connection.SendMessage, Data: repositories.Message{Message: "Hi"}}
		return req
	}

	return mockConn1, mockConn2
}

func TestManager_NotifyInMemory(t *testing.T) {
	td := newTestData()
	messageRepo := td.MockMessageRepo
	convoRepo := td.MockConversationRepo

	messageRepo.Create = func(message repositories.Message) (message2 *repositories.Message, e error) {
		return &message, nil
	}

	convoRepo.GetConvo = func(conversationId string) (conversants []repositories.Conversant, e error) {
		return []repositories.Conversant{
			{ID: "a"},
		}, nil
	}

	manager := NewManager(td, td, td, td)

	sender := make(chan struct{})
	conn1, conn2 := inMemoryMockCons(sender)

	manager.addConn(conn1)
	manager.addConn(conn2)

	if len(manager.connections) != 2 {
		t.Error("didn't add connections")
	}

	sentChannel := make(chan connection.Response, 1)

	conn2.Resp = func() chan connection.Response {
		return sentChannel
	}

	sender <- struct{}{}

	select {
	case <-sentChannel:
		t.Log("passed")
	case <-time.After(1 * time.Second):
		t.Error("Message not recieved")
	}
}

func TestManager_Leave(t *testing.T) {
	td := newTestData()
	manager := NewManager(td, td, td, td)

	conn1 := makeConn("a")

	closer := make(chan struct{})

	conn1.Leaver = func() chan struct{} {
		return closer
	}

	manager.addConn(conn1)

	if len(manager.connections) == 0 {
		t.Error("conn didn't join")
	}

	closer <- struct{}{}

	if len(manager.connections) != 0 {
		t.Error("conn didn't leave")
	}
}
