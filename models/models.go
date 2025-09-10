package models

// Message представляет серверное сообщение, отправляемое клиентам по WS
type Message struct {
	Nickname string `json:"nickname"`
	Text     string `json:"text"`
	Time     string `json:"time"`

	// Для изображений
	ImageURL string `json:"imageUrl,omitempty"`

	// Для произвольных файлов
	FileURL  string `json:"fileUrl,omitempty"`
	FileName string `json:"fileName,omitempty"`
	FileSize int64  `json:"fileSize,omitempty"`
	MimeType string `json:"mimeType,omitempty"`

	// "text" | "image" | "file"
	Type string `json:"type"`
}

// ClientMessage представляет сообщение, приходящее от клиента по WS
type ClientMessage struct {
	Text string `json:"text"`

	// Для изображений
	ImageURL string `json:"imageUrl,omitempty"`

	// Для произвольных файлов
	FileURL  string `json:"fileUrl,omitempty"`
	FileName string `json:"fileName,omitempty"`
	FileSize int64  `json:"fileSize,omitempty"`
	MimeType string `json:"mimeType,omitempty"`

	// "text" | "image" | "file"
	Type string `json:"type"`
}
