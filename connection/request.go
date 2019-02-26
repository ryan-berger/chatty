package connection

type RequestType int
type ResponseType int

const (
	SendMessage RequestType = iota
	CreateConversation
)

const (
	Error ResponseType = iota
	NewMessage
	NewConversation
	MessageSent
)

type Request struct {
	Type RequestType
	Data interface{}
}

type Response struct {
	Type ResponseType
	Data interface{}
}

type ResponseError struct {
	Error string
}

func NewResponseError(error string) Response {
	return Response{
		Type: Error,
		Data: ResponseError{Error: error},
	}
}
