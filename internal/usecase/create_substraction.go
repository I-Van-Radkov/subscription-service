package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/dto"
	"github.com/I-Van-Radkov/subscription-service/internal/models"
	"github.com/google/uuid"
)

func (u *SubscriptionUsecase) CreateSubscription(ctx context.Context, input dto.CreateSubstractionRequest) (dto.CreateSubstractionResponse, error) {
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return dto.CreateSubstractionResponse{}, fmt.Errorf("invalid user_id format: %w", err)
	}

	start, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return dto.CreateSubstractionResponse{}, fmt.Errorf("invalid start_date format (expected MM-YYYY): %w", err)
	}

	var end *time.Time
	if input.EndDate != nil {
		t, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return dto.CreateSubstractionResponse{}, fmt.Errorf("invalid end_date format (expected MM-YYYY): %w", err)
		}
		end = &t
	}

	sub := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     end,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := sub.Validate(); err != nil {
		return dto.CreateSubstractionResponse{}, fmt.Errorf("validation model failed: %w", err)
	}

	id, err := u.Repository.Create(ctx, sub)
	if err != nil {
		return dto.CreateSubstractionResponse{}, fmt.Errorf("db failed to create subscription: %w", err)
	}

	output := dto.CreateSubstractionResponse{
		ID: id.String(),
	}

	return output, nil
}
