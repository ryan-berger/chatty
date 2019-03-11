package repositories

// Conversant is a struct representing someone who converses
// A conversant is not unique per connection, but is distinct
// from a user as a Conversant is not organizationally dependent
type Conversant struct {
	ID          string `json:"id"db:"id"`
	DisplayName string `json:"name"db:"display_name"`
}

// Message is an incoming message to be sent to all conversants
// within the given conversation
type Message struct {
	ID             string `json:"id"db:"id"`
	SenderID       string `json:"senderId"db:"sender_id"`
	Message        string `json:"message"db:"message"`
	ConversationID string `json:"conversationId"db:"conversation_id"`
}

// Conversation is a group of conversants, and a list of messages
// allowing the manager to make sure everyone is notified of a message
type Conversation struct {
	ID          string       `json:"id"db:"id"`
	Name        string       `json:"id"db:"name"`
	Conversants []Conversant `json:"conversants"`
	Messages    []Message    `json:"messages"`
	Direct      bool         `json:"direct"db:"direct"`
}
