package connection

import (
	"github.com/ryan-berger/chatty/repositories"
)

// Conn is the generic connection interface that allows
// multiple connections to talk to each other over any
// protocol
type Conn interface {
	Authorize() error
	GetConversant() repositories.Conversant
	Requests() chan Request
	Response() chan Response
	Leave() chan struct{}
}

type MockConn struct {
	Auth       func() error
	Conversant func() repositories.Conversant
	Request    func() chan Request
	Resp       func() chan Response
	Leaver     func() chan struct{}
}

func (mock *MockConn) Authorize() error {
	return mock.Auth()
}

func (mock *MockConn) Requests() chan Request {
	return mock.Request()
}

func (mock *MockConn) Response() chan Response {
	return mock.Resp()
}

func (mock *MockConn) Leave() chan struct{} {
	return mock.Leaver()
}

func (mock *MockConn) GetConversant() repositories.Conversant {
	return mock.Conversant()
}
