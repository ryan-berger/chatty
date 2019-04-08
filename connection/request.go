package connection

import (
	"errors"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type RequestType int
type ResponseType int

const (
	SendMessage RequestType = iota
	CreateConversation
	RetrieveConversation
	RequestError
)

type (
	// Request is a struct with a dynamic body, with a specific type
	// that allows the MUX to figure out how to cast the body
	Request struct {
		Type RequestType `json:"type"`
		Data interface{} `json:"data"`
	}

	// CreateConversationRequest takes in a name and list of extra users
	// to add to the conversation. SenderID is for internal use only,
	// and should not be marshalled
	CreateConversationRequest struct {
		SenderID    string   `json:"-"`
		Name        string   `json:"name"`
		Conversants []string `json:"conversants"`
	}

	// SendMessageRequest takes the message and conversation
	// ID and sends the given message to the conversation
	SendMessageRequest struct {
		SenderID       string `json:"-"`
		Message        string `json:"message"`
		ConversationID string `json:"conversationId"`
	}

	// RetrieveConversationRequest uses a limit offset pattern in order to return
	// chats from a conversation
	RetrieveConversationRequest struct {
		ConversationID string `json:"conversationId"`
		Limit          int    `json:"limit"`
		Offset         int    `json:"offset"`
	}
)

// UUIDList valdiates
func UUIDList(input interface{}) error {
	if list, ok := input.([]string); ok {
		for _, s := range list {
			err := is.UUIDv4.Validate(s)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("must be list type")
	}
	return nil
}

func (request CreateConversationRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Conversants, validation.Length(2, 20), validation.By(UUIDList)))
}

func (request SendMessageRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Message, validation.Required),
		validation.Field(&request.ConversationID, is.UUIDv4),
		validation.Field(&request.ConversationID, validation.Required, is.UUIDv4))
}

func (request RetrieveConversationRequest) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.ConversationID, validation.Required, is.UUIDv4),
		validation.Field(&request.Limit, validation.Required, validation.Min(1)),
		validation.Field(&request.Offset, validation.Required, validation.Min(0)))
}
