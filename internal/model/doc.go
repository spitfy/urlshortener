// doc.go
package model

// URL документация для Swagger
type URL struct {
	// Оригинальный URL
	OriginalURL string `json:"original_url"`
	// Сокращенный URL
	ShortURL string `json:"short_url"`
}
