package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/dto"
	"github.com/google/uuid"
)

func (u *SubscriptionUsecase) UpdateSubscription(ctx context.Context, idString string, input dto.UpdateSubscriptionRequest) (dto.UpdateSubscriptionResponse, error) {
	idUUID, err := uuid.Parse(idString)
	if err != nil {
		return dto.UpdateSubscriptionResponse{}, fmt.Errorf("invalid id format: %w", err)
	}

	// Проверка на существование подписки
	sub, err := u.Repository.GetById(ctx, idUUID)
	if err != nil {
		return dto.UpdateSubscriptionResponse{}, err
	}
	if sub == nil {
		return dto.UpdateSubscriptionResponse{}, fmt.Errorf("subscription not found")
	}

	if input.ServiceName != "" {
		sub.ServiceName = input.ServiceName
	}
	if input.Price > 0 {
		sub.Price = input.Price
	}
	if input.StartDate != "" {
		start, err := time.Parse("01-2006", input.StartDate)
		if err != nil {
			return dto.UpdateSubscriptionResponse{}, fmt.Errorf("invalid start_date: %w", err)
		}
		sub.StartDate = start
	}
	if input.EndDate != nil {
		if *input.EndDate == "" {
			sub.EndDate = nil
		} else {
			t, err := time.Parse("01-2006", *input.EndDate)
			if err != nil {
				return dto.UpdateSubscriptionResponse{}, fmt.Errorf("invalid end_date: %w", err)
			}
			sub.EndDate = &t
		}
	}

	sub.UpdatedAt = time.Now()

	if err := u.Repository.Update(ctx, sub); err != nil {
		return dto.UpdateSubscriptionResponse{}, fmt.Errorf("failed update subscription: %w", err)
	}

	var endDate *string
	if sub.EndDate != nil {
		t := sub.EndDate.Format("01-2006")
		endDate = &t
	}

	output := dto.UpdateSubscriptionResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     endDate,
	}

	return output, nil
}
