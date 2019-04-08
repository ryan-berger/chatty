package connection

import "github.com/ryan-berger/chatty/repositories"

const (
	Error ResponseType = iota
	NewMessage
	NewConversation
	ReturnConversation
)

type (
	Response struct {
		Type ResponseType `json:"type"`
		Data interface{}  `json:"data"`
	}

	NewMessageResponse struct {
		repositories.Message
	}

	ResponseError struct {
		Error string `json:"error"`
	}
)

func NewResponseError(error string) Response {
	return Response{
		Type: Error,
		Data: ResponseError{Error: error},
	}
}
