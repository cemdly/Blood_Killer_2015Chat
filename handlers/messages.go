package handlers

import (
	"chat-room/backend/config"
	"chat-room/backend/models"
	"chat-room/backend/storage"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func sendInitialData(ws *websocket.Conn, nickname string) error {
	if err := sendHistory(ws); err != nil {
		return err
	}

	handleMessage("Система", models.ClientMessage{
		Text: nickname + " подключился к чату",
		Type: "text",
	})

	return nil
}

func processMessages(ws *websocket.Conn, nickname string) {
	for {
		msg, err := readClientMessage(ws)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		if err := handleMessage(nickname, msg); err != nil {
			log.Printf("Message handling error: %v", err)
		}
	}
}

func readClientMessage(ws *websocket.Conn) (models.ClientMessage, error) {
	var msg models.ClientMessage
	err := ws.ReadJSON(&msg)
	return msg, err
}

func handleMessage(nickname string, msg models.ClientMessage) error {
	var err error

	switch msg.Type {
	case "image":
		if msg.ImageURL != "" {
			err = storage.StorageRepo.SaveImageMessage(nickname, msg.Text, msg.ImageURL)
		} else {
			err = storage.StorageRepo.SaveMessage(nickname, msg.Text)
		}
	case "file":
		if msg.FileURL != "" {
			err = storage.StorageRepo.SaveFileMessage(
				nickname,
				msg.Text,
				msg.FileURL,
				msg.FileName,
				msg.FileSize,
				msg.MimeType,
			)
		} else {
			err = storage.StorageRepo.SaveMessage(nickname, msg.Text)
		}
	default:
		err = storage.StorageRepo.SaveMessage(nickname, msg.Text)
	}

	if err != nil {
		return err
	}

	broadcastMsg := models.Message{
		Nickname: nickname,
		Text:     msg.Text,
		Time:     time.Now().Format(time.RFC3339),
		Type:     msg.Type,
	}

	if msg.Type == "image" {
		broadcastMsg.ImageURL = msg.ImageURL
	}
	if msg.Type == "file" {
		broadcastMsg.FileURL = msg.FileURL
		broadcastMsg.FileName = msg.FileName
		broadcastMsg.FileSize = msg.FileSize
		broadcastMsg.MimeType = msg.MimeType
	}

	config.Broadcast <- broadcastMsg
	return nil
}

func sendHistory(ws *websocket.Conn) error {
	messages, err := storage.StorageRepo.GetLastMessages(50)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err := ws.WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func MessagesHandler() {
	for msg := range config.Broadcast {
		config.ClientsMu.Lock()
		for client := range config.Clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Write error: %v", err)
				client.Close()
				delete(config.Clients, client)
			}
		}
		config.ClientsMu.Unlock()
	}
}
