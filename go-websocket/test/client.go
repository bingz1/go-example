package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:7777", "http service address")

func loop() {

	for {
		u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			continue
		}
		// 循环读消息
		for {
			c.WriteMessage(websocket.TextMessage, []byte("hello"))
			_, message, err := c.ReadMessage()
			if err != nil {
				// log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
		}
		c.Close()
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	for i := 0; i < 1000; i++ {
		go loop()
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
