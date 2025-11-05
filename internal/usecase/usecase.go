package usecase

import (
	"context"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/models"
	"github.com/google/uuid"
)

type SubscriptionRepo interface {
	Create(ctx context.Context, sub *models.Subscription) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error)
	SumForPeriod(ctx context.Context, userID uuid.UUID, serviceName string, start time.Time, end *time.Time) (int, error)
}

type SubscriptionUsecase struct {
	Repository SubscriptionRepo
}

func NewSubscriptionUsecase(repo SubscriptionRepo) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		Repository: repo,
	}
}
