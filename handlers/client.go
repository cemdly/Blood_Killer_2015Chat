package handlers

import (
	"chat-room/backend/config"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func getNickname(r *http.Request) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := make([]byte, 3)
	for i := range digits {
		digits[i] = byte(rand.Intn(10)) + '0'
	}

	if nickname := r.URL.Query().Get("nickname"); nickname != "" {
		if !isNickTaken(nickname) {
			return nickname
		} else {
			return nickname + string(digits)
		}
	}
	return getRandNickname()
}

func getRandNickname() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := make([]byte, 3)
	for i := range digits {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	randNick := "Kidder" + string(digits)

	return randNick
}

func isNickTaken(value string) bool {
	for _, v := range config.Clients {
		if v == value {
			return true
		}
	}
	return false
}

func registerClient(ws *websocket.Conn, nickname string) {
	config.ClientsMu.Lock()
	config.Clients[ws] = nickname
	config.ClientsMu.Unlock()
	log.Printf("%s connected!", nickname)
}

func unregisterClient(ws *websocket.Conn) {
	config.ClientsMu.Lock()
	nickname := config.Clients[ws]
	delete(config.Clients, ws)
	config.ClientsMu.Unlock()
	log.Printf("%s disconnected!", nickname)
}
