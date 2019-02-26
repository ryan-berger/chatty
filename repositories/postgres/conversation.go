package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/ryan-berger/chatty/repositories"
	"github.com/ryan-berger/chatty/repositories/models"
)

const CREATECONVERSATION = `
INSERT INTO conversation(id, name, direct) VALUES (:id, :name, :direct)
`

const CREATECONVERSANTCONVERSATION = `
INSERT INTO conversant_conversation(conversation_id, conversant_id) VALUES ($1, $2)
`

const GETUSERSFROMCONVESRATION = `
SELECT
    id,
	external_id
FROM conversant_conversation cc
LEFT JOIN conversant c on cc.conversant_id = c.id
WHERE conversation_id = $1
`

const GETCONVERSATIONMESSAGES = `
SELECT
	m.id,
    m.sender,
    m.message
FROM chat_message m
WHERE m.conversation = $1
LIMIT $2 OFFSET $3
`

type ConversationRepository struct {
	db *sqlx.DB
}

func (repo *ConversationRepository) CreateConversation(conversation models.Conversation) (*models.Conversation, error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}

	conversation.Id = uuid.New()
	conversation.Direct = len(conversation.Conversants) == 2
	_, err = tx.NamedExec(CREATECONVERSATION, &conversation)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "err: creating conversation")
	}

	err = addConversants(tx, conversation.Id, conversation.Conversants)
	if err != nil {
		return nil, errors.Wrap(err, "err: Adding Conversants")
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "commit: error comitting CreateConversation")
	}

	return &conversation, nil
}

func addConversants(tx *sqlx.Tx, conversationId string, conversants []models.Conversant) error {
	for _, conversant := range conversants {
		_, err := tx.Exec(CREATECONVERSANTCONVERSATION, &conversationId, &conversant.Id)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (repo *ConversationRepository) RetrieveConversation(conversationId string, limit, offset int) (*models.Conversation, error) {
	var conversation models.Conversation

	err := repo.db.Get(&conversation, CREATECONVERSANTCONVERSATION)

	if err != nil {
		return nil, err
	}

	err = repo.db.Select(&(conversation.Conversants), GETUSERSFROMCONVESRATION, &conversationId)
	if err != nil {
		return nil, err
	}

	err = repo.db.Get(&(conversation.Messages), GETCONVERSATIONMESSAGES, &conversationId, &limit, &offset)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (repo *ConversationRepository) GetConversants(conversationId string) ([]models.Conversant, error)  {
	var conversants []models.Conversant
	err := repo.db.Select(&conversants, GETUSERSFROMCONVESRATION, &conversationId)
	if err != nil {
		return nil, err
	}

	return conversants, nil
}

func NewConversationRepository(db *sqlx.DB) repositories.ConversationRepo  {
	return &ConversationRepository{
		db:db,
	}
}
