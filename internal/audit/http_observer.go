package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// HTTPObserver отправляет события аудита на указанный HTTP endpoint
type HTTPObserver struct {
	url string
}

// NewHTTPObserver создает новый HTTPObserver с настраиваемым HTTP клиентом
func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{url: url}
}

// Notify отправляет событие аудита на HTTP endpoint
func (o *HTTPObserver) Notify(ctx context.Context, event Event) error {
	if o.url == "" {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return nil
}
