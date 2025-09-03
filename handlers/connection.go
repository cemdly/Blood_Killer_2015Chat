package handlers

import (
	"chat-room/backend/config"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection(w, r)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	nickname := getNickname(r)
	registerClient(ws, nickname)
	defer unregisterClient(ws)

	if err := sendInitialData(ws, nickname); err != nil {
		log.Printf("Initial data error: %v", err)
	}

	processMessages(ws, nickname)
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return config.Upgrader.Upgrade(w, r, nil)
}
