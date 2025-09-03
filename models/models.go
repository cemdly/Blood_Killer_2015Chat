package models

type Message struct {
	Nickname string `json:"nickname"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	ImageURL string `json:"imageUrl,omitempty"`
	Type     string `json:"type"` // "text" или "image"
}

type ClientMessage struct {
	Text     string `json:"text"`
	ImageURL string `json:"imageUrl,omitempty"`
	Type     string `json:"type"` // "text" или "image"
}
