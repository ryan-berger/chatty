package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/ryan-berger/chatty/repositories"
)

const createConversation = `
INSERT INTO conversation(id, name, direct) VALUES (:id, :name, :direct)
`

const createConversantConversation = `
INSERT INTO conversant_conversation(conversation_id, conversant_id) VALUES ($1, $2)
`

const getUsersFromConversation = `
SELECT
    id,
    name
FROM conversant_conversation cc
LEFT JOIN conversant c on cc.conversant_id = c.id
WHERE conversation_id = $1
`

const getConversationMessages = `
SELECT
	m.id,
    m.sender,
    m.message
FROM chat_message m
WHERE m.conversation = $1
LIMIT $2 OFFSET $3
`

// ConversationRepository is an implementation of ConversationRepo
// that uses Postgres as it's backend
type ConversationRepository struct {
	db *sqlx.DB
}

// CreateConversation creates a conversation with a postgres transaction. If any of it fails, it
// rolls back or returns an error
func (repo *ConversationRepository) CreateConversation(conversation repositories.Conversation) (*repositories.Conversation, error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}

	conversation.ID = uuid.New()
	conversation.Direct = len(conversation.Conversants) == 2
	_, err = tx.NamedExec(createConversation, &conversation)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "err: creating conversation")
	}

	err = addConversants(tx, conversation.ID, conversation.Conversants)
	if err != nil {
		return nil, errors.Wrap(err, "err: Adding Conversants")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "commit: error comitting CreateConversation")
	}

	return &conversation, nil
}

func addConversants(tx *sqlx.Tx, conversationID string, conversants []repositories.Conversant) error {
	for _, conversant := range conversants {
		_, err := tx.Exec(createConversantConversation, &conversationID, &conversant.ID)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

// RetrieveConversation grabs a conversation with the messages given a limit and offset
func (repo *ConversationRepository) RetrieveConversation(conversationID string, limit, offset int) (*repositories.Conversation, error) {
	var conversation repositories.Conversation

	err := repo.db.Get(&conversation, createConversantConversation)

	if err != nil {
		return nil, err
	}

	err = repo.db.Select(&(conversation.Conversants), getUsersFromConversation, &conversationID)
	if err != nil {
		return nil, err
	}

	err = repo.db.Get(&(conversation.Messages), getConversationMessages, &conversationID, &limit, &offset)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

// GetConversants gets the conversants for a given conversation
func (repo *ConversationRepository) GetConversants(conversationID string) ([]repositories.Conversant, error) {
	var conversants []repositories.Conversant
	err := repo.db.Select(&conversants, getUsersFromConversation, &conversationID)
	if err != nil {
		return nil, err
	}

	return conversants, nil
}

// NewConversationRepository creates a Postgres instance of a ConversationRepo
func NewConversationRepository(db *sqlx.DB) *ConversationRepository {
	return &ConversationRepository{
		db: db,
	}
}
