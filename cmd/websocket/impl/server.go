package impl

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ryan-berger/chatty/manager"
	"log"
	"net/http"
	"net/http/pprof"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunServer(m *manager.Manager) {
	http.HandleFunc("/", pprof.Index)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(m, writer, request)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func serveWs(manager *manager.Manager, writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)

	fmt.Println("upgrade")

	if err != nil {
		log.Println("err on upgrade", err)
		return
	}

	manager.Join(newWsConn(conn))
}


