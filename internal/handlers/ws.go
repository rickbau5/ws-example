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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		// ensure connection is closed by forcibly closing it
		if err := conn.Close(); err != nil {
			handler.Logger.Println("Error closing connection:", err)
		}
	}()

	err = conn.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf("[%s] hello :) type a message and press ENTER", handler.ID)),
	)

	if err != nil {
		handler.Logger.Println("Error writing greeting to socket:", err)
		// continue on regardless of error
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			handler.Logger.Println("Error reading message:", err)
			return
		}
		handler.Logger.Println("Got message:", strings.Trim(string(p), "\n"))
		err = conn.WriteMessage(messageType,
			append([]byte(fmt.Sprintf("[%s] ", handler.ID)), p...),
		)
		if err != nil {
			handler.Logger.Println("Error writing message to socket:", err)
			return
		}
	}
}
