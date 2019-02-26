package noop

import (
	"fmt"

	"github.com/ryan-berger/chatty/repositories"
)

type Notifier struct {
}

func (*Notifier) Notify(id string, message repositories.Message) error {
	fmt.Printf("To: %s Message: %s\n", id, message.Message)
	return nil
}

func NewNotifier() *Notifier {
	return &Notifier{}
}
