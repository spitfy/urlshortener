package audit

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

func ExampleFileObserver_Notify() {
	// Создаем временный файл для логов
	tmpFile, err := os.CreateTemp("", "audit_log_example.txt")
	if err != nil {
		fmt.Println("Failed to create temp file:", err)
		return
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	// Создаем FileObserver
	observer := NewFileObserver(tmpFile.Name())

	// Создаем событие аудита
	event := Event{
		Timestamp: time.Now(),
		Action:    Shorten,
		UserID:    123,
		URL:       "https://example.com",
	}

	// Отправляем событие
	err = observer.Notify(context.Background(), event)
	if err != nil {
		fmt.Println("Error notifying observer:", err)
		return
	}

	// Читаем содержимое файла
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println("Failed to read temp file:", err)
		return
	}

	// Формируем и выводим только часть строки для проверки
	line := strings.TrimSpace(string(content))
	parts := strings.SplitN(line, " ", 4)
	if len(parts) == 4 {
		fmt.Printf("%s %s %s\n", parts[1], parts[2], parts[3])
	} else {
		fmt.Println(line)
	}

	// Output: shorten 123 https://example.com
}
