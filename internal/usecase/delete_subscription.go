package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *SubscriptionUsecase) DeleteSubscription(ctx context.Context, idString string) error {
	idUUID, err := uuid.Parse(idString)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	if err := u.Repository.Delete(ctx, idUUID); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
