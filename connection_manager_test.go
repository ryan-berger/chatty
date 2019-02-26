package chatty

import (
	"testing"
	"time"

	"github.com/ryan-berger/chatty/manager/connection"
	"github.com/ryan-berger/chatty/manager/operators"
	"github.com/ryan-berger/chatty/repositories"
	"github.com/ryan-berger/chatty/repositories/models"
)

type testData struct {
	*repositories.MockConversationRepo
	*repositories.MockMessageRepo
	*connection.MockAuther
	*operators.MockNotifier
}

func newTestData() *testData {
	return &testData{
		&repositories.MockConversationRepo{},
		&repositories.MockMessageRepo{},
		&connection.MockAuther{},
		&operators.MockNotifier{},
	}
}

func makeConn(id string) *connection.MockConn {
	mockConn := &connection.MockConn{}
	mockConn.Conversant = func() models.Conversant {
		return models.Conversant{
			Id: id,
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
		req <- connection.Request{Type: connection.SendMessage, Data: models.Message{Message: "Hi"}}
		return req
	}

	return mockConn1, mockConn2
}

func TestManager_NotifyInMemory(t *testing.T) {
	td := newTestData()
	messageRepo := td.MockMessageRepo
	convoRepo := td.MockConversationRepo

	messageRepo.Create = func(message models.Message) (message2 *models.Message, e error) {
		return &message, nil
	}

	convoRepo.GetConvo = func(conversationId string) (conversants []models.Conversant, e error) {
		return []models.Conversant{
			{Id: "a"},
		}, nil
	}

	manager := NewManager(
		repositories.NewChatInteractor(td, td), td, td)

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
	manager := NewManager(repositories.NewChatInteractor(td, td), td, td)

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
