package repositories

import "github.com/ryan-berger/chatty/repositories/models"

type MessageRepo interface {
	CreateMessage(message models.Message) (*models.Message, error)
}

type MockMessageRepo struct {
	Create func(message models.Message) (*models.Message, error)
}

func (mock *MockMessageRepo) CreateMessage(message models.Message) (*models.Message, error) {
	return mock.Create(message)
}
