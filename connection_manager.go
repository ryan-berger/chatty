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
	data connection.SendMessageRequest
}

// ConnectionManager is the main connection manager struct that handles all chat connections
type ConnectionManager struct {
	auther         operators.Auther
	connectionMu   *sync.RWMutex
	connections    map[string][]connection.Conn
	messageChan    chan messageRequest
	shutdownChan   chan struct{}
	chatInteractor *chatInteractor
	notifier       operators.Notifier
}

// NewManager creates a new connection manager given repos and operators
func NewManager(
	messageRepo repositories.MessageRepo,
	conversationRepo repositories.ConversationRepo,
	conversantRepo repositories.ConversantRepo,
	auther operators.Auther,
	notifier operators.Notifier) *ConnectionManager {

	manager := &ConnectionManager{
		auther:         auther,
		connectionMu:   &sync.RWMutex{},
		connections:    make(map[string][]connection.Conn),
		shutdownChan:   make(chan struct{}),
		messageChan:    make(chan messageRequest, numWorkers),
		chatInteractor: newChatInteractor(messageRepo, conversationRepo, conversantRepo),
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

	_, err := manager.chatInteractor.UpsertConvserant(conn.GetConversant())

	if err != nil {
		fmt.Println(err)
		conn.Leave() <- struct{}{}
		return
	}

	manager.addConn(conn)
}

func (manager *ConnectionManager) addConn(conn connection.Conn) {
	fmt.Println("joining")
	manager.connectionMu.Lock()
	connections := manager.connections[conn.GetConversant().ID]
	manager.connections[conn.GetConversant().ID] = append(connections, conn)
	go manager.handleConnection(conn)
	manager.connectionMu.Unlock()
}

func (manager *ConnectionManager) handleConnection(conn connection.Conn) {
	for {
		select {
		case command := <-conn.Requests():
			if command.Data == nil {
				manager.sendErr(conn, "no request body")
				continue
			}
			switch command.Type {
			case connection.SendMessage:
				manager.sendMessage(conn, command.Data.(connection.SendMessageRequest))
			case connection.CreateConversation:
				manager.createConversation(conn, command.Data.(connection.CreateConversationRequest))
			}
		case <-conn.Leave():
			manager.removeConn(conn.GetConversant().ID)
			return
		}
	}
}

func (manager *ConnectionManager) removeConn(id string) {
	manager.connectionMu.Lock()
	defer manager.connectionMu.Unlock()

	for i, clientConn := range manager.connections[id] {
		connArray := manager.connections[id]
		if id == clientConn.GetConversant().ID {
			fmt.Println("before: ", connArray)
			connArray[i] = connArray[len(connArray)-1]
			fmt.Println("after: ", connArray[:len(connArray)-1])
			manager.connections[id] = connArray[:len(connArray)-1]

		}
	}

	if len(manager.connections[id]) == 0 {
		delete(manager.connections, id)
	}
}

func (manager *ConnectionManager) sendMessage(conn connection.Conn, m connection.SendMessageRequest) {
	for {
		select {
		case manager.messageChan <- messageRequest{conn: conn, data: m}:
			return
		default:
			continue
		}
	}
}

func (manager *ConnectionManager) createConversation(sender connection.Conn, conversation connection.CreateConversationRequest) {
	newConversation, err := manager.
		chatInteractor.
		CreateConversation(sender.GetConversant().ID, conversation.Name, conversation.Conversants)

	if err != nil {
		manager.sendErr(sender, "unable to create conversation")
	}

	sender.Response() <- connection.Response{Type: connection.NewConversation, Data: *newConversation}
}

func (manager *ConnectionManager) startMessageWorker() {
	for {
		select {
		case message := <-manager.messageChan:
			manager.createMessage(message)
		case <-manager.shutdownChan:
			return
		}
	}
}

func (manager *ConnectionManager) createMessage(message messageRequest) {
	messageData := message.data
	newMessage, err := manager.
		chatInteractor.
		SendMessage(messageData.Message, message.conn.GetConversant().ID, messageData.ConversationID)

	if err != nil {
		manager.sendErr(message.conn, "couldn't send message")
		return
	}

	conversants, err := manager.chatInteractor.GetConversants(messageData.ConversationID)

	if err != nil {
		return
	}

	manager.notifyRecipients(conversants, *newMessage)
}

func (manager *ConnectionManager) notifyRecipients(conversants []repositories.Conversant, message repositories.Message) {
	manager.connectionMu.RLock()
	for _, conversant := range conversants {
		if val, ok := manager.connections[conversant.ID]; ok {
			for _, conn := range val {
				manager.sendNewMessage(conn, message)
			}
		} else {
			manager.notifier.Notify(conversant.ID, message)
		}
	}
	manager.connectionMu.RUnlock()
}

func (manager *ConnectionManager) sendNewMessage(conn connection.Conn, message repositories.Message) {
	conn.Response() <- connection.Response{Type: connection.NewMessage, Data: message}
}

func (manager *ConnectionManager) sendErr(conn connection.Conn, errString string) {
	manager.connectionMu.RLock()
	conn.Response() <- connection.NewResponseError(errString)
	manager.connectionMu.RUnlock()
}
