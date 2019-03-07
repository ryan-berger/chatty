package repositories

type ConversantRepo interface {
	UpdateOrCreate(conversant Conversant) (*Conversant, error)
}
