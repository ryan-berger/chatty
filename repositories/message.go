package repositories

// MessageRepo is a way for the connection manager to store messages
type MessageRepo interface {
	CreateMessage(message Message) (*Message, error)
}

// MockMessageRepo is a MessageRepo implementation for testing
type MockMessageRepo struct {
	Create func(message Message) (*Message, error)
}

// CreateMessage calls the Create method in the MockMessageRepo
func (mock *MockMessageRepo) CreateMessage(message Message) (*Message, error) {
	return mock.Create(message)
}
