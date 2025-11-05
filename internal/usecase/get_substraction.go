package usecase

import (
	"context"
	"fmt"

	"github.com/I-Van-Radkov/subscription-service/internal/dto"
	"github.com/google/uuid"
)

func (u *SubscriptionUsecase) GetSubscription(ctx context.Context, idString string) (dto.GetSubscriptionResponse, error) {
	idUUID, err := uuid.Parse(idString)
	if err != nil {
		return dto.GetSubscriptionResponse{}, fmt.Errorf("invalid id format: %w", err)
	}

	sub, err := u.Repository.GetById(ctx, idUUID)
	if err != nil {
		return dto.GetSubscriptionResponse{}, fmt.Errorf("failed to get subscription: %w", err)
	}
	if sub == nil {
		return dto.GetSubscriptionResponse{}, fmt.Errorf("subscription not found")
	}

	var endDate *string
	if sub.EndDate != nil {
		t := sub.EndDate.Format("01-2006")
		endDate = &t
	}

	output := dto.GetSubscriptionResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     endDate,
	}

	return output, nil
}
