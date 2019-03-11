package connection

type RequestType int
type ResponseType int

const (
	SendMessage RequestType = iota
	CreateConversation
	RequestError
)

type (
	Request struct {
		Type RequestType `json:"type"`
		Data interface{} `json:"data"`
	}

	CreateConversationRequest struct {
		Name        string   `json:"name"`
		Conversants []string `json:"conversants"`
	}

	SendMessageRequest struct {
		Message        string `json:"message"`
		ConversationID string `json:"conversationId"`
	}
)
