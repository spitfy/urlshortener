package audit

import (
	"context"
	"os"
	"sync"
	"time"
)

type FileObserver struct {
	filePath string
	mu       sync.Mutex
}

func NewFileObserver(filePath string) *FileObserver {
	return &FileObserver{filePath: filePath}
}

func (o *FileObserver) Notify(ctx context.Context, event Event) error {
	if o.filePath == "" {
		return nil
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	file, err := os.OpenFile(o.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(event.Timestamp.Format(time.RFC3339) + " " + event.Method + " " + event.Hash + " " + event.Link + "\n")
	return err
}
