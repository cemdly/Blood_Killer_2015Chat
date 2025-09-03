package config

import (
	"log"
	"net/http"
	"os"
)

var BldPath = "/home/student210/Документы/Zhuravlev22-40/docs/Go-WebChat/frontend/build"

func ServeFrontend() {
	if _, err := os.Stat(BldPath); os.IsNotExist(err) {
		log.Fatal("React build directory not found. Please run 'npm run build' in frontend directory first.")
	}

	fs := http.FileServer(http.Dir(BldPath))
	http.Handle("/", fs)

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, BldPath+"/index.html")
	})
}

// func GetPort() string {
// 	if port := os.Getenv("PORT"); port != "" {
// 		return port
// 	}
// 	return "8080"
// }
