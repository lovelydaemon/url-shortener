package usecase

import (
	"context"
	"fmt"
)

type PingUseCase struct {
	repo PingRepo
}

func NewPingUseCase(r PingRepo) *PingUseCase {
	return &PingUseCase{r}
}

func (uc *PingUseCase) Ping(ctx context.Context) error {
	if err := uc.repo.Ping(ctx); err != nil {
		return fmt.Errorf("PingUseCase - Ping - uc.repo.Ping: %w", err)
	}
	return nil
}
