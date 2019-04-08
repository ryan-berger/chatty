package chatty

import (
	"github.com/pborman/uuid"
	"github.com/ryan-berger/chatty/connection"

	"github.com/ryan-berger/chatty/repositories"
)

type chatInteractor struct {
	conversationRepo repositories.ConversationRepo
	messageRepo      repositories.MessageRepo
	conversantRepo   repositories.ConversantRepo
}

func (chat *chatInteractor) CreateConversation(request connection.CreateConversationRequest) (*repositories.Conversation, error) {
	request.Conversants = append(request.Conversants, request.SenderID)

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	newConversation := repositories.Conversation{}
	for _, conversantID := range request.Conversants {
		newConversation.Conversants = append(newConversation.Conversants, repositories.Conversant{ID: conversantID})
	}

	newConversation.Name = request.Name
	newConversation.Direct = len(newConversation.Conversants) == 2

	convo, err := chat.conversationRepo.CreateConversation(newConversation)
	if err != nil {
		return nil, err
	}

	return convo, nil
}

func (chat *chatInteractor) GetConversation(request connection.RetrieveConversationRequest) (*repositories.Conversation, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	conversation, err := chat.conversationRepo.RetrieveConversation(request.ConversationID, request.Limit, request.Offset)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (chat *chatInteractor) SendMessage(message connection.SendMessageRequest) (*repositories.Message, error) {

	err := message.Validate()

	if err != nil {
		return nil, err
	}

	msg := repositories.Message{
		ID:             uuid.New(),
		Message:        message.Message,
		SenderID:       message.SenderID,
		ConversationID: message.ConversationID,
	}

	newMessage, err := chat.messageRepo.CreateMessage(msg)

	if err != nil {
		return nil, err
	}

	return newMessage, nil
}

func (chat *chatInteractor) GetConversants(conversationID string) ([]repositories.Conversant, error) {
	conversants, err := chat.conversationRepo.GetConversants(conversationID)

	if err != nil {
		return nil, err
	}

	return conversants, nil
}

func (chat *chatInteractor) UpsertConvserant(conversant repositories.Conversant) (*repositories.Conversant, error) {
	newConversant, err := chat.conversantRepo.UpdateOrCreate(conversant)

	if err != nil {
		return nil, err
	}

	return newConversant, nil
}

func newChatInteractor(
	messageRepo repositories.MessageRepo,
	conversationRepo repositories.ConversationRepo,
	conversantRepo repositories.ConversantRepo) *chatInteractor {
	return &chatInteractor{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		conversantRepo:   conversantRepo,
	}
}
