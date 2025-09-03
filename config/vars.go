package config

import (
	"chat-room/backend/models"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Clients   = make(map[*websocket.Conn]string)
	ClientsMu sync.Mutex
	Broadcast = make(chan models.Message)
	Upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
