package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/ryan-berger/chatty/repositories/models"
)

type MessageRepository struct {
	db *sqlx.DB
}

func (*MessageRepository) CreateMessage(message models.Message) (*models.Message, error) {
	panic("implement me")
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository  {
	return &MessageRepository{
		db:db,
	}
}