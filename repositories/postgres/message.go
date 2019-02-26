package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/ryan-berger/chatty/repositories"
)

// MessageRepository is a MessageRepo implementation that uses Postgres to store messages
type MessageRepository struct {
	db *sqlx.DB
}

// CreateMessage stores a message in Postgres
func (*MessageRepository) CreateMessage(message repositories.Message) (*repositories.Message, error) {
	panic("implement me")
}

// NewMessageRepository creates a new Postgres MessageRepository
func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}
