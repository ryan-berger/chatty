package repositories

import (
	"errors"
	"github.com/ryan-berger/chatty/repositories/models"
)

type ChatInteractor struct {
	conversationRepo ConversationRepo
	messageRepo      MessageRepo
}

var INTERNAL = errors.New("err: internal server error")

func (chat *ChatInteractor) CreateConversation(conversation models.Conversation) (*models.Conversation, error) {
	convo, err := chat.conversationRepo.CreateConversation(conversation)
	if err != nil {
		return nil, err
	}

	return convo, nil
}

func (chat *ChatInteractor) GetConversation(id string, offset, limit int) (*models.Conversation, error) {
	conversation, err := chat.conversationRepo.RetrieveConversation(id, offset, limit)

	if err != nil {
		return nil, INTERNAL
	}

	return conversation, nil
}

func (chat *ChatInteractor) SendMessage(message models.Message) (*models.Message, error) {
	newMessage, err := chat.messageRepo.CreateMessage(message)

	if err != nil {
		return nil, INTERNAL
	}

	return newMessage, nil
}

func (chat *ChatInteractor) GetConversants(conversationId string) ([]models.Conversant, error) {
	conversants, err := chat.conversationRepo.GetConversants(conversationId)

	if err != nil {
		return nil, INTERNAL
	}

	return conversants, nil
}

func NewChatInteractor(messageRepo MessageRepo, conversationRepo ConversationRepo) *ChatInteractor {
	return &ChatInteractor{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}
