package chatty

import (
	"github.com/ryan-berger/chatty/interactors"
	"github.com/ryan-berger/chatty/repositories"
)

type Configuration struct {
	Auther         interactors.Auther
	ChatInteractor *interactors.Chat
}

func NewConfiguration(
	auther interactors.Auther,
	convoRepo repositories.ConversationRepo,
	messageRepo repositories.MessageRepo, ) *Configuration {
	return &Configuration{
		Auther: auther,
		ChatInteractor: interactors.NewChatInteractor(auther, messageRepo, convoRepo),
	}
}
