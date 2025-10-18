package audit

import (
	"context"
	"os"
	"sync"
	"time"
)

// FileObserver реализует Observer для записи событий в файл
type FileObserver struct {
	filePath string
	mu       sync.Mutex
}

// NewFileObserver создает новый FileObserver
// filePath: путь к файлу журнала. Если пустой - логирование отключено.
func NewFileObserver(filePath string) *FileObserver {
	return &FileObserver{filePath: filePath}
}

// Notify записывает событие в файл журнала
func (o *FileObserver) Notify(_ context.Context, event Event) error {
	if o.filePath == "" {
		return nil
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	file, err := os.OpenFile(o.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = file.WriteString(event.Timestamp.Format(time.RFC3339) + " " + string(event.Action) + " " + string(rune(event.UserID)) + " " + event.URL + "\n")
	return err
}
