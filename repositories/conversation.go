package repositories

// ConversationRepo is a way for the connection manager to store conversations
type ConversationRepo interface {
	CreateConversation(conversation Conversation) (*Conversation, error)
	RetrieveConversation(conversationID string, limit, offset int) (*Conversation, error)
	GetConversants(conversationID string) ([]Conversant, error)
}

// MockConversationRepo is a mock conversation repo for testing
type MockConversationRepo struct {
	CreateConvo   func(conversation Conversation) (*Conversation, error)
	RetrieveConvo func(conversationId string, limit, offset int) (*Conversation, error)
	GetConvo      func(conversationId string) ([]Conversant, error)
}

// CreateConversation calls CreateConvo inside of the MockConversationRepo struct
func (m *MockConversationRepo) CreateConversation(conversation Conversation) (*Conversation, error) {
	return m.CreateConvo(conversation)
}

// RetrieveConversation calls RetrieveConvo in the MockConversationRepo struct
func (m *MockConversationRepo) RetrieveConversation(conversationID string, limit, offset int) (*Conversation, error) {
	return m.RetrieveConvo(conversationID, limit, offset)
}

// GetConversants calls GetConvo in the MockConversationRepo struct
func (m *MockConversationRepo) GetConversants(conversationID string) ([]Conversant, error) {
	return m.GetConvo(conversationID)
}
