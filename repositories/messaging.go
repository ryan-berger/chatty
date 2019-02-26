package repositories

// Conversant is a struct representing someone who converses
// A conversant is not unique per connection, but is distinct
// from a user as a Conversant is not organizationally dependent
type Conversant struct {
	ID   string
	Name string
}

// Message is an incoming message to be sent to all conversants
// within the given conversation
type Message struct {
	ID             string
	SenderID       string
	Message        string
	ConversationID string
}

// Conversation is a group of conversants, and a list of messages
// allowing the manager to make sure everyone is notified of a message
type Conversation struct {
	ID          string
	Conversants []Conversant
	Messages    []Message
	Direct      bool
}
