package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/dto"
	"github.com/google/uuid"
)

func (u *SubscriptionUsecase) GetSubscriptionsSum(ctx context.Context, userIdStr, serviceName, start string, end *string) (dto.GetSubSumResponse, error) {
	var userId uuid.UUID
	if userIdStr != "" {
		parsed, err := uuid.Parse(userIdStr)
		if err != nil {
			return dto.GetSubSumResponse{}, fmt.Errorf("invalid user_id: %w", err)
		}
		userId = parsed
	}

	startDate, err := time.Parse("01-2006", start)
	if err != nil {
		return dto.GetSubSumResponse{}, fmt.Errorf("invalid start_date format (expected MM-YYYY): %w", err)
	}

	var endDate *time.Time
	if end != nil {
		t, err := time.Parse("01-2006", *end)
		if err != nil {
			return dto.GetSubSumResponse{}, fmt.Errorf("invalid end_date format (expected MM-YYYY): %w", err)
		}
		endDate = &t
	}

	total, err := u.Repository.SumForPeriod(ctx, userId, serviceName, startDate, endDate)
	if err != nil {
		return dto.GetSubSumResponse{}, fmt.Errorf("failed to get summary of period from DB: %w", err)
	}

	output := dto.GetSubSumResponse{
		Total: total,
	}

	return output, nil
}
