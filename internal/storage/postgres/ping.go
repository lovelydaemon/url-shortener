package postgres

import (
	"context"
	"fmt"
)

func (s storage) Ping(ctx context.Context) error {
	if err := s.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("Storage - Ping - s.Pool.Ping: %w", err)
	}
	return nil
}
