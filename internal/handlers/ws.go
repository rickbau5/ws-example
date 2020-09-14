package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	ID       string
	Upgrader websocket.Upgrader
	Logger   *log.Logger
}

func (handler *WebSocketHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := handler.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		handler.Logger.Println("Error upgrading connection:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf("[%s] hello :) type a message and press ENTER", handler.ID)),
	)

	if err != nil {
		log.Println(err)
		// continue on regardless of error
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Got message:", strings.Trim(string(p), "\n"))
		err = conn.WriteMessage(messageType,
			append([]byte(fmt.Sprintf("[%s] ", handler.ID)), p...),
		)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
