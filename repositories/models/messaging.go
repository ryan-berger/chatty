package models

type Conversant struct {
	Id   string
	Name string
}

type Message struct {
	Id             string
	SenderId       string
	Message        string
	ConversationId string
}

type Conversation struct {
	Id          string
	Conversants []Conversant
	Messages    []Message
	Direct      bool
}
