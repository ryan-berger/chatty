package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ryan-berger/chatty/cmd/websocket/impl"
	"github.com/ryan-berger/chatty/manager"
	"github.com/ryan-berger/chatty/manager/operators/noop"
	"github.com/ryan-berger/chatty/repositories"
	"github.com/ryan-berger/chatty/repositories/postgres"
	"os"
)

func GetDBString() string {
	return fmt.Sprintf(
		"host=%s database=%s user=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))
}

func main() {
	db, err := sqlx.Open("postgres", GetDBString())

	if err != nil {
		panic(err)
	}

	conversationRepo := postgres.NewConversationRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	chatInteractor := repositories.NewChatInteractor(messageRepo, conversationRepo)
	notifier := noop.NewNotifier()
	man := manager.NewManager(chatInteractor, nil, notifier)

	impl.RunServer(man)
}
