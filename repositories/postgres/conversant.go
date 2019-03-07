package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/ryan-berger/chatty/repositories"
)

const updateOrCreateConversant = `
INSERT INTO conversant (id, "name") VALUES (:id, :name) 
ON CONFLICT (id) DO 
  UPDATE SET  = :name
`

type ConversantRepository struct {
	db *sqlx.DB
}

func (repo *ConversantRepository) UpdateOrCreate(conversant repositories.Conversant) (*repositories.Conversant, error) {
	_, err := repo.db.Exec(updateOrCreateConversant)

	if err != nil {
		return nil, err
	}

	return &conversant, nil
}
