package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/ryan-berger/chatty/repositories"
)

const createMessage = `
INSERT INTO  chat_message(id, message, sender, conversation) 
VALUES (:id, :message, :sender, :conversation)
`

// MessageRepository is a MessageRepo implementation that uses Postgres to store messages
type MessageRepository struct {
	db *sqlx.DB
}

// CreateMessage stores a message in Postgres
func (repo *MessageRepository) CreateMessage(message repositories.Message) (*repositories.Message, error) {
	message.ID = uuid.New()
	_, err := repo.db.NamedExec(createMessage, &message)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// NewMessageRepository creates a new Postgres MessageRepository
func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}
