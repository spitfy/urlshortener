package audit

import (
	"context"
	"time"
)

type Event struct {
	Timestamp time.Time
	Method    string
	Hash      string
	Link      string
}

type Observer interface {
	Notify(ctx context.Context, event Event) error
}
