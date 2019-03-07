package chatty

import (
	"errors"

	"github.com/ryan-berger/chatty/repositories"
)

type chatInteractor struct {
	conversationRepo repositories.ConversationRepo
	messageRepo      repositories.MessageRepo
	conversantRepo   repositories.ConversantRepo
}

// Signifying an internal server error
var errInternal = errors.New("err: internal server error")

func (chat *chatInteractor) CreateConversation(conversation repositories.Conversation) (*repositories.Conversation, error) {
	convo, err := chat.conversationRepo.CreateConversation(conversation)
	if err != nil {
		return nil, err
	}

	return convo, nil
}

func (chat *chatInteractor) GetConversation(id string, offset, limit int) (*repositories.Conversation, error) {
	conversation, err := chat.conversationRepo.RetrieveConversation(id, offset, limit)

	if err != nil {
		return nil, errInternal
	}

	return conversation, nil
}

func (chat *chatInteractor) SendMessage(message repositories.Message) (*repositories.Message, error) {
	newMessage, err := chat.messageRepo.CreateMessage(message)

	if err != nil {
		return nil, errInternal
	}

	return newMessage, nil
}

func (chat *chatInteractor) GetConversants(conversationID string) ([]repositories.Conversant, error) {
	conversants, err := chat.conversationRepo.GetConversants(conversationID)

	if err != nil {
		return nil, errInternal
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
