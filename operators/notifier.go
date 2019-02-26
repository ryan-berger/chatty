package operators

import "github.com/ryan-berger/chatty/repositories"

// Notifier interface that notifies a user when they
// are not logged on
type Notifier interface {
	Notify(id string, message repositories.Message) error
}

// MockNotifier operator for testing
type MockNotifier struct {
	SendNotification func(id string, message repositories.Message) error
}

// Notify using SendNotificaion in MockNotifier struct
func (mock *MockNotifier) Notify(id string, message repositories.Message) error {
	return mock.SendNotification(id, message)
}
