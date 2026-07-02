package audit

import "context"

type Event struct {
	TenantID   string
	ActorID    string
	Action     string
	EntityType string
	EntityID   string
}

type Recorder interface {
	Record(ctx context.Context, event Event) error
}
