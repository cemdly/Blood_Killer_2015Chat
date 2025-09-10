package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	UploadDir   = "uploads"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func HandleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверяем размер файла
	if header.Size > MaxFileSize {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Проверяем тип файла
	if !isValidImageType(header) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Генерируем уникальное имя файла
	fileName, err := generateUniqueFileName(header.Filename)
	if err != nil {
		log.Printf("Error generating filename: %v", err)
		http.Error(w, "Error generating filename", http.StatusInternalServerError)
		return
	}

	// Создаем файл на диске
	filePath := filepath.Join(UploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое файла
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error copying file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Возвращаем URL файла
	imageURL := fmt.Sprintf("/uploads/%s", fileName)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"imageUrl": "%s"}`, imageURL)
}

func isValidImageType(header *multipart.FileHeader) bool {
	contentType := header.Header.Get("Content-Type")
	return allowedImageTypes[contentType]
}

func generateUniqueFileName(originalName string) (string, error) {
	// Генерируем случайные байты
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Получаем расширение файла
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}

	// Создаем уникальное имя
	randomString := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%s%s", randomString, ext), nil
}

func ServeUploadedFiles(w http.ResponseWriter, r *http.Request) {
	// Убираем префикс /uploads/ из пути
	filePath := strings.TrimPrefix(r.URL.Path, "/uploads/")
	
	// Проверяем, что путь не содержит .. (защита от directory traversal)
	if strings.Contains(filePath, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(UploadDir, filePath)
	
	// Проверяем, что файл существует
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Определяем content-type по расширению
	ext := strings.ToLower(filepath.Ext(filePath))
	ctype := mime.TypeByExtension(ext)
	if ctype == "" {
		ctype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", ctype)

	http.ServeFile(w, r, fullPath)
}

// HandleFileUpload — универсальная загрузка любых файлов (поле формы: "file")
func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(MaxFileSize); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	fileName, err := generateUniqueFileName(header.Filename)
	if err != nil {
		log.Printf("Error generating filename: %v", err)
		http.Error(w, "Error generating filename", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(UploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Определяем mime-тип по расширению
	ext := strings.ToLower(filepath.Ext(header.Filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	fileURL := fmt.Sprintf("/uploads/%s", fileName)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"fileUrl":"%s","fileName":"%s","fileSize":%d,"mimeType":"%s"}`,
		fileURL, header.Filename, written, mimeType)
}
