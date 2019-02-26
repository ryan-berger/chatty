package operators

import "github.com/ryan-berger/chatty/repositories/models"

type Notifier interface {
	Notify(id string, message models.Message) error
}

type MockNotifier struct {
	SendNotification func(id string, message models.Message) error
}

func (mock *MockNotifier) Notify(id string, message models.Message) error {
	return mock.SendNotification(id, message)
}
