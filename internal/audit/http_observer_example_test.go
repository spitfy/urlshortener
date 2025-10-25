package audit

import (
	"context"
	"log"
	"time"
)

func ExampleHTTPObserver_Notify() {
	// Создаем наблюдателя с тестовым url
	observer := NewHTTPObserver("http://example.com/audit")

	// Формируем тестовое событие
	event := Event{
		Timestamp: time.Now(),
		Action:    Shorten,
		UserID:    123,
		URL:       "http://short.url/abc",
	}

	// Отправляем событие
	err := observer.Notify(context.Background(), event)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
}
