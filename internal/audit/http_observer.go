package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type HTTPObserver struct {
	url string
}

func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{url: url}
}

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

	_, err = http.DefaultClient.Do(req)
	return err
}
