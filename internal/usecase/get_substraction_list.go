package usecase

import (
	"context"
	"fmt"

	"github.com/I-Van-Radkov/subscription-service/internal/dto"
	"github.com/google/uuid"
)

func (u *SubscriptionUsecase) GetSubscriptionsList(ctx context.Context, userIDStr string) (dto.GetSubsListResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.GetSubsListResponse{}, fmt.Errorf("invalid user_id: %w", err)
	}

	subs, err := u.Repository.List(ctx, userID)
	if err != nil {
		return dto.GetSubsListResponse{}, err
	}

	output := dto.GetSubsListResponse{
		Total: 0,
		List:  make([]dto.GetSubscriptionResponse, 0, len(subs)),
	}
	for _, sub := range subs {
		var end *string
		if sub.EndDate != nil {
			t := sub.EndDate.Format("01-2006")
			end = &t
		}

		output.List = append(output.List, dto.GetSubscriptionResponse{
			ID:          sub.ID.String(),
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID.String(),
			StartDate:   sub.StartDate.Format("01-2006"),
			EndDate:     end,
		})
		output.Total++
	}

	return output, nil
}
