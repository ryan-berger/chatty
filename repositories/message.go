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

// DefaultMockRepo creates a mock repo that will return
// any data given to it without errors
func DefaultMockMessageRepo() MessageRepo {
	return &MockMessageRepo{
		Create: func(message Message) (*Message, error) {
			return &message, nil
		},
	}
}
