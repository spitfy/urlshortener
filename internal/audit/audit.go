package audit

import (
	"context"
	"time"
)

// Action представляет тип действия пользователя для аудита
type Action string

const (
	Shorten Action = "shorten" // Действие: сокращение URL
	Follow  Action = "follow"  // Действие: переход по сокращенному URL
)

// Event содержит информацию о событии для аудита
type Event struct {
	Timestamp time.Time `json:"ts"`      // Временная метка события
	Action    Action    `json:"action"`  // Тип действия
	UserID    int       `json:"user_id"` // ID пользователя
	URL       string    `json:"url"`     // URL, к которому относится действие
}

// Observer определяет интерфейс для наблюдателей аудита
type Observer interface {
	// Notify отправляет событие аудита
	// Возвращает ошибку, если отправка не удалась
	Notify(ctx context.Context, event Event) error
}
