package connection

import "github.com/ryan-berger/chatty/repositories"

const (
	Error ResponseType = iota
	NewMessage
	NewConversation
	MessageSent
)

type (
	Response struct {
		Type ResponseType
		Data interface{}
	}

	NewMessageResponse struct {
		repositories.Message
	}

	ResponseError struct {
		Error string
	}
)

func NewResponseError(error string) Response {
	return Response{
		Type: Error,
		Data: ResponseError{Error: error},
	}
}
