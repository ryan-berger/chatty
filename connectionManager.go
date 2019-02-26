package chatty

import (
	"fmt"
	"sync"

	"github.com/ryan-berger/chatty/repositories"

	"github.com/ryan-berger/chatty/connection"
	"github.com/ryan-berger/chatty/operators"
)

var numWorkers = 40

type messageRequest struct {
	conn connection.Conn
	data repositories.Message
}

// ConnectionManager is the main connection manager struct that handles all chat connections
type ConnectionManager struct {
	auther         operators.Auther
	connectionMu   *sync.RWMutex
	connections    map[string]connection.Conn
	messageChan    chan messageRequest
	shutdownChan   chan struct{}
	chatInteractor *chatInteractor
	notifier       operators.Notifier
}

// NewManager creates a new connection manager given repos and operators
func NewManager(
	messageRepo repositories.MessageRepo,
	conversationRepo repositories.ConversationRepo,
	auther operators.Auther,
	notifier operators.Notifier) *ConnectionManager {

	manager := &ConnectionManager{
		auther:         auther,
		connectionMu:   &sync.RWMutex{},
		connections:    make(map[string]connection.Conn),
		shutdownChan:   make(chan struct{}),
		messageChan:    make(chan messageRequest, numWorkers),
		chatInteractor: &chatInteractor{messageRepo: messageRepo, conversationRepo: conversationRepo},
		notifier:       notifier,
	}
	manager.startup()
	return manager
}

func (manager *ConnectionManager) startup() {
	for i := 0; i < numWorkers; i++ {
		go manager.startMessageWorker()
	}
}

func (manager *ConnectionManager) shutdown() {
	manager.shutdownChan <- struct{}{}
}

// Join authorizes a connection and then joins the server
func (manager *ConnectionManager) Join(conn connection.Conn) {
	if err := conn.Authorize(); err != nil {
		return
	}

	manager.addConn(conn)
}

func (manager *ConnectionManager) addConn(conn connection.Conn) {
	fmt.Println("joining")
	manager.connectionMu.Lock()
	manager.connections[conn.GetConversant().ID] = conn
	go manager.handleConnection(conn)
	manager.connectionMu.Unlock()
}

func (manager *ConnectionManager) closeConn(conn connection.Conn) {
	conn.Leave() <- struct{}{}
}

func (manager *ConnectionManager) handleConnection(conn connection.Conn) {
	for {
		select {
		case command := <-conn.Requests():
			switch command.Type {
			case connection.SendMessage:
				manager.sendMessage(conn, command.Data.(repositories.Message))
			case connection.CreateConversation:
				manager.createConversation(command.Data.(repositories.Conversation))
			}
		case <-conn.Leave():
			manager.connectionMu.Lock()
			delete(manager.connections, conn.GetConversant().ID)
			manager.connectionMu.Unlock()
			return
		}
	}
}

func (manager *ConnectionManager) sendMessage(conn connection.Conn, m repositories.Message) {
	for {
		select {
		case manager.messageChan <- messageRequest{conn: conn, data: m}:
			return
		default:
			continue
		}
	}
}

func (manager *ConnectionManager) createConversation(conversation repositories.Conversation) {

}

func (manager *ConnectionManager) startMessageWorker() {
	for {
		select {
		case message := <-manager.messageChan:
			if message.conn.GetConversant().ID != message.data.SenderID {
				manager.sendErr(message.conn, "can't send message for someone else")
			}

			_, err := manager.
				chatInteractor.
				SendMessage(message.data)

			if err != nil {
				manager.sendErr(message.conn, "couldn't send message")
				continue
			}

			conversants, err := manager.chatInteractor.GetConversants(message.data.ConversationID)

			if err != nil {
				continue
			}

			manager.notifyRecipients(conversants, message.data)
		case <-manager.shutdownChan:
			return
		default:
			continue
		}
	}
}

func (manager *ConnectionManager) notifyRecipients(conversants []repositories.Conversant, message repositories.Message) {
	manager.connectionMu.RLock()
	for _, conversant := range conversants {
		if val, ok := manager.connections[conversant.ID]; ok {
			manager.sendNewMessage(val, message)
		} else {
			manager.notifier.Notify(conversant.ID, message)
		}
	}
	manager.connectionMu.RUnlock()
}

func (manager *ConnectionManager) sendNewMessage(conn connection.Conn, message repositories.Message) {

}

func (manager *ConnectionManager) sendErr(conn connection.Conn, errString string) {
	manager.connectionMu.RLock()
	conn.Response() <- connection.NewResponseError(errString)
	manager.connectionMu.RUnlock()
}
