package storage

import (
	"chat-room/backend/models"
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var StorageRepo *Storage

type Storage struct {
	db *sql.DB
}

func InitStorage(store *Storage) {
	StorageRepo = store
}

func New(connStr string) (*Storage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id SERIAL PRIMARY KEY,
            nickname TEXT NOT NULL,
            text TEXT NOT NULL,
            image_url TEXT,
            message_type TEXT NOT NULL DEFAULT 'text',
            file_url TEXT,
            file_name TEXT,
            file_size BIGINT,
            mime_type TEXT,
            timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )`)
	return err
}

func (s *Storage) SaveMessage(nickname, text string) error {
	_, err := s.db.Exec(
		"INSERT INTO messages (nickname, text, message_type) VALUES ($1, $2, 'text')",
		nickname, text)
	return err
}

func (s *Storage) SaveImageMessage(nickname, text, imageURL string) error {
	_, err := s.db.Exec(
		"INSERT INTO messages (nickname, text, image_url, message_type) VALUES ($1, $2, $3, 'image')",
		nickname, text, imageURL)
	return err
}

func (s *Storage) SaveFileMessage(nickname, text, fileURL, fileName string, fileSize int64, mimeType string) error {
	_, err := s.db.Exec(
		"INSERT INTO messages (nickname, text, file_url, file_name, file_size, mime_type, message_type) VALUES ($1, $2, $3, $4, $5, $6, 'file')",
		nickname, text, fileURL, fileName, fileSize, mimeType)
	return err
}

func (s *Storage) GetLastMessages(limit int) ([]models.Message, error) {
	rows, err := s.db.Query(`
        SELECT 
            nickname, 
            text, 
            image_url,
            message_type,
            file_url,
            file_name,
            file_size,
            mime_type,
            TO_CHAR(timestamp AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS time
        FROM messages 
        ORDER BY timestamp DESC 
        LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		var imageURL sql.NullString
		var fileURL sql.NullString
		var fileName sql.NullString
		var fileSize sql.NullInt64
		var mimeType sql.NullString
		if err := rows.Scan(&m.Nickname, &m.Text, &imageURL, &m.Type, &fileURL, &fileName, &fileSize, &mimeType, &m.Time); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		if imageURL.Valid {
			m.ImageURL = imageURL.String
		}
		if fileURL.Valid {
			m.FileURL = fileURL.String
		}
		if fileName.Valid {
			m.FileName = fileName.String
		}
		if fileSize.Valid {
			m.FileSize = fileSize.Int64
		}
		if mimeType.Valid {
			m.MimeType = mimeType.String
		}
		messages = append(messages, m)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
