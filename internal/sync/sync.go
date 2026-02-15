package sync

import "context"

type Syncer interface {
	Sync(ctx context.Context) error
}

type NoOpSyncer struct{}

func NewNoOpSyncer() *NoOpSyncer {
	return &NoOpSyncer{}
}

func (s *NoOpSyncer) Sync(ctx context.Context) error {
	// TODO: implement real sync
	return nil
}
