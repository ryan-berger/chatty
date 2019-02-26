package connection

import (
	"github.com/ryan-berger/chatty/repositories/models"
)

type Auther interface {
	CanStartConversation(conversation models.Conversation) bool
}

type MockAuther struct {
	CanStartConvo func(models.Conversation) bool
}

func (mock *MockAuther) CanStartConversation(conversation models.Conversation) bool {
	return mock.CanStartConvo(conversation)
}
