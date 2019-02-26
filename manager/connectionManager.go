package manager

import (
	"fmt"
	"github.com/ryan-berger/chatty/manager/connection"
	"github.com/ryan-berger/chatty/manager/operators"
	"github.com/ryan-berger/chatty/repositories"
	"github.com/ryan-berger/chatty/repositories/models"
	"sync"
)

const NUMWORKERS = 40

type messageRequest struct {
	conn connection.Conn
	data models.Message
}

type Manager struct {
	auther         connection.Auther
	connectionMu   *sync.RWMutex
	connections    map[string]connection.Conn
	messageChan    chan messageRequest
	shutdownChan   chan struct{}
	chatInteractor *repositories.ChatInteractor
	notifier       operators.Notifier
}

func NewManager(chat *repositories.ChatInteractor, auther connection.Auther, notifier operators.Notifier) *Manager {
	manager := &Manager{
		auther:         auther,
		connectionMu:   &sync.RWMutex{},
		connections:    make(map[string]connection.Conn),
		shutdownChan:   make(chan struct{}),
		messageChan:    make(chan messageRequest, NUMWORKERS),
		chatInteractor: chat,
		notifier:       notifier,
	}
	manager.startup()
	return manager
}

func (manager *Manager) startup() {
	for i := 0; i < NUMWORKERS; i++ {
		go manager.startMessageWorker()
	}
}

func (manager *Manager) shutdown() {
	manager.shutdownChan <- struct{}{}
}

func (manager *Manager) Join(conn connection.Conn) {
	if err := conn.Authorize(); err != nil {
		return
	}

	manager.addConn(conn)
}


func (manager *Manager) addConn(conn connection.Conn) {
	fmt.Println("joining")
	manager.connectionMu.Lock()
	manager.connections[conn.GetConversant().Id] = conn
	go manager.handleConnection(conn)
	manager.connectionMu.Unlock()
}

func (manager *Manager) closeConn(conn connection.Conn) {
	conn.Leave() <- struct{}{}
}

func (manager *Manager) handleConnection(conn connection.Conn) {
	for {
		select {
		case command := <-conn.Requests():
			switch command.Type {
			case connection.SendMessage:
				manager.sendMessage(conn, command.Data.(models.Message))
			case connection.CreateConversation:
				manager.createConversation(command.Data.(models.Conversation))
			}
		case <-conn.Leave():
			manager.connectionMu.Lock()
			delete(manager.connections, conn.GetConversant().Id)
			manager.connectionMu.Unlock()
			return
		}
	}
}

func (manager *Manager) sendMessage(conn connection.Conn, m models.Message) {
	for {
		select {
		case manager.messageChan <- messageRequest{conn: conn, data: m}:
			return
		default:
			continue
		}
	}
}

func (manager *Manager) createConversation(conversation models.Conversation) {

}

func (manager *Manager) startMessageWorker() {
	for {
		select {
		case message := <-manager.messageChan:
			if message.conn.GetConversant().Id != message.data.SenderId {
				manager.sendErr(message.conn, "can't send message for someone else")
			}

			_, err := manager.
				chatInteractor.
				SendMessage(message.data)

			if err != nil {
				manager.sendErr(message.conn, "couldn't send message")
				continue
			}

			conversants, err := manager.chatInteractor.GetConversants(message.data.ConversationId)

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

func (manager *Manager) notifyRecipients(conversants []models.Conversant, message models.Message) {
	manager.connectionMu.RLock()
	for _, conversant := range conversants {
		if val, ok := manager.connections[conversant.Id]; ok {
			manager.sendNewMessage(val, message)
		} else {
			manager.notifier.Notify(conversant.Id, message)
		}
	}
	manager.connectionMu.RUnlock()
}

func (manager *Manager) sendNewMessage(conn connection.Conn, message models.Message) {

}

func (manager *Manager) sendErr(conn connection.Conn, errString string) {
	manager.connectionMu.RLock()
	conn.Response() <- connection.NewResponseError(errString)
	manager.connectionMu.RUnlock()
}
