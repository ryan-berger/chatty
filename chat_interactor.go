package chatty

import (
	"github.com/pborman/uuid"

	"github.com/ryan-berger/chatty/repositories"
)

type chatInteractor struct {
	conversationRepo repositories.ConversationRepo
	messageRepo      repositories.MessageRepo
	conversantRepo   repositories.ConversantRepo
}

func (chat *chatInteractor) CreateConversation(creatorID string, name string, conversants []string) (*repositories.Conversation, error) {
	newConversation := repositories.Conversation{}

	conversants = append(conversants, creatorID)

	for _, conversantID := range conversants {
		newConversation.Conversants = append(newConversation.Conversants, repositories.Conversant{ID: conversantID})
	}
	newConversation.Name = name

	convo, err := chat.conversationRepo.CreateConversation(newConversation)
	if err != nil {
		return nil, err
	}

	return convo, nil
}

func (chat *chatInteractor) GetConversation(id string, offset, limit int) (*repositories.Conversation, error) {
	conversation, err := chat.conversationRepo.RetrieveConversation(id, offset, limit)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (chat *chatInteractor) SendMessage(message string, sender string, conversationId string) (*repositories.Message, error) {

	msg := repositories.Message{
		ID:             uuid.New(),
		Message:        message,
		SenderID:       sender,
		ConversationID: conversationId,
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
