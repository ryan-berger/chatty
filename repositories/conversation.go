package repositories

import "github.com/ryan-berger/chatty/repositories/models"

type ConversationRepo interface {
	CreateConversation(conversation models.Conversation) (*models.Conversation, error)
	RetrieveConversation(conversationId string, limit, offset int) (*models.Conversation, error)
	GetConversants(conversationId string) ([]models.Conversant, error)
}

type MockConversationRepo struct {
	CreateConvo   func(conversation models.Conversation) (*models.Conversation, error)
	RetrieveConvo func(conversationId string, limit, offset int) (*models.Conversation, error)
	GetConvo      func(conversationId string) ([]models.Conversant, error)
}

func (m *MockConversationRepo) CreateConversation(conversation models.Conversation) (*models.Conversation, error) {
	return m.CreateConvo(conversation)
}

func (m *MockConversationRepo) RetrieveConversation(conversationId string, limit, offset int) (*models.Conversation, error) {
	return m.RetrieveConvo(conversationId, limit, offset)
}

func (m *MockConversationRepo) GetConversants(conversationId string) ([]models.Conversant, error) {
	return m.GetConvo(conversationId)
}
