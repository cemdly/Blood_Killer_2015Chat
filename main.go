package main

import (
	"chat-room/backend/config"
	"chat-room/backend/handlers"
	"chat-room/backend/storage"
	"log"
	"net/http"
)

func main() {

	DbUrl := "postgresql://postgres:277353@localhost:5432/chatdb"

	store, err := storage.New(DbUrl)
	if err != nil {
		log.Fatal("Failed to initialize storage: ", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("Error closing storage: %v", err)
		}
	}()

	storage.InitStorage(store)

	go handlers.MessagesHandler()

	http.HandleFunc("/ws", handlers.HandleConnections)
	http.HandleFunc("/upload", handlers.HandleImageUpload)
	http.HandleFunc("/uploads/", handlers.ServeUploadedFiles)

	config.ServeFrontend()

	//port := config.GetPort()
	log.Printf("Server started")
	log.Fatal(http.ListenAndServe("0.0.0.0:2773", nil))
}
