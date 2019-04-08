package repositories

type ConversantRepo interface {
	UpdateOrCreate(conversant Conversant) (*Conversant, error)
}

type MockConversantRepo struct {
	Upsert func(conversant Conversant) (*Conversant, error)
}

func (repo *MockConversantRepo) UpdateOrCreate(conversant Conversant) (*Conversant, error) {
	return repo.Upsert(conversant)
}
