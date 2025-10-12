package audit

import (
	"context"
	"time"
)

type Action string

const (
	Shorten Action = "shorten"
	Follow  Action = "follow"
)

type Event struct {
	Timestamp time.Time `json:"ts"`
	Action    Action    `json:"action"`
	UserID    int       `json:"user_id"`
	URL       string    `json:"url"`
}

type Observer interface {
	Notify(ctx context.Context, event Event) error
}
