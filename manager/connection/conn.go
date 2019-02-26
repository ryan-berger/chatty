package connection

import "github.com/ryan-berger/chatty/repositories/models"

type Conn interface {
	Authorize() error
	GetConversant() models.Conversant
	Requests() chan Request
	Response() chan Response
	Leave() chan struct{}
}

type MockConn struct {
	Auth       func() error
	Conversant func() models.Conversant
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

func (mock *MockConn) GetConversant() models.Conversant {
	return mock.Conversant()
}
