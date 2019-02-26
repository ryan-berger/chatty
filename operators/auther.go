package operators

import (
	"github.com/ryan-berger/chatty/repositories"
)

// Auther Interface helps figure out if a conversation can be started
type Auther interface {
	CanStartConversation(conversation repositories.Conversation) bool
}

// MockAuther operator for testing
type MockAuther struct {
	CanStartConvo func(repositories.Conversation) bool
}

// CanStartConversation calls function passed in from MockAuther for testing
func (mock *MockAuther) CanStartConversation(conversation repositories.Conversation) bool {
	return mock.CanStartConvo(conversation)
}
