package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			conn, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost:8080")
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}

			fmt.Println(conn)

			var b []byte

			conn.Read(b)
			conn.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}
