package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gorilla/websocket"
	ws "github.com/ryan-berger/chatty/connection/websocket"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ryan-berger/chatty"
	"github.com/ryan-berger/chatty/operators/noop"
	"github.com/ryan-berger/chatty/repositories/postgres"
)

func getDBString() string {
	return fmt.Sprintf(
		"host=%s database=%s user=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	db, err := sqlx.Open("postgres", getDBString())

	if err != nil {
		panic(err)
	}

	conversationRepo := postgres.NewConversationRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	notifier := noop.NewNotifier()
	man := chatty.NewManager(messageRepo, conversationRepo, nil, notifier)

	http.HandleFunc("/", pprof.Index)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(man, writer, request)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveWs(manager *chatty.Manager, writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)

	fmt.Println("upgrade")

	if err != nil {
		log.Println("err on upgrade", err)
		return
	}

	manager.Join(ws.NewWebsocketConn(conn))
}
