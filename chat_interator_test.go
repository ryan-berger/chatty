package chatty

import (
	"fmt"
	"testing"

	"github.com/ryan-berger/chatty/repositories"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/ryan-berger/chatty/connection"
)

var formError = map[string]map[string]string{
	"AllErrors": {
		"name":        "cannot be blank",
		"conversants": "the length must be between 2 and 20",
	},
	"UUIDList": {
		"conversants": "must be a valid UUID v4",
	},
}

func validateForm(t *testing.T, err validation.Errors, errorSet string) {
	for k, v := range formError[errorSet] {
		if val, ok := err[k]; !ok {
			t.Fatalf("%s: expected %s, received none", k, v)
		} else {
			if v != val.Error() {
				t.Fatalf("%s: expected %s, received: %s", k, v, val)
			}
		}
	}
}

func TestChatInteractor_CreateConversationValidation(t *testing.T) {
	interactor := &chatInteractor{}
	request := connection.CreateConversationRequest{
		SenderID: "d8ece527-a0e9-4513-8972-5a7b0f97785d",
	}
	_, err := interactor.CreateConversation(request)
	validateForm(t, err.(validation.Errors), "AllErrors")

	request.SenderID = "asdf"
	request.Conversants = []string{"d8ece527-a0e9-4513-8972-5a7b0f97785d"}
	_, err = interactor.CreateConversation(request)
	validateForm(t, err.(validation.Errors), "UUIDList")
}

func TestChatInteractor_CreateConversation(t *testing.T) {
	convoRepo := &repositories.MockConversationRepo{
		CreateConvo: func(conversation repositories.Conversation) (conversation2 *repositories.Conversation, e error) {
			return &conversation, nil
		},
	}
	interactor := &chatInteractor{
		conversationRepo: convoRepo,
	}

	request := connection.CreateConversationRequest{
		SenderID:    "d8ece527-a0e9-4513-8972-5a7b0f97785d",
		Name:        "Test",
		Conversants: []string{"d8ece527-a0e9-4513-8972-5a7b0f97785d"},
	}

	convo, err := interactor.CreateConversation(request)

	if err != nil {
		fmt.Println(err)
		t.Fatalf("Create convo shouldn't have failed")
	}

	if len(convo.Conversants) != 2 {
		t.Fatalf("Second conversant not appended")
	}

	if convo.Direct != true {
		t.Fatalf("Conversation should be direct")
	}
}

func TestChatInteractor_GetConversation(t *testing.T) {
	called := false
	convoRepo := &repositories.MockConversationRepo{
		RetrieveConvo: func(conversationId string, limit, offset int) (conversation *repositories.Conversation, e error) {
			called = true
			return &repositories.Conversation{
				ID: "a",
			}, nil
		},
	}

	interactor := &chatInteractor{
		conversationRepo: convoRepo,
	}

	interactor.GetConversation(connection.RetrieveConversationRequest{})

	if called {
		t.Fatalf("Retrieve Conversation shouldn't have been called")
	}
}
