package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (s *Subscription) Validate() error {
	if s.ServiceName == "" {
		return errors.New("service_name is required")
	}

	if s.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if s.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if s.StartDate.IsZero() {
		return errors.New("start_date is required")
	}

	if s.EndDate != nil {
		if s.EndDate.Before(s.StartDate) {
			return errors.New("end_date cannot be before start_date")
		}
	}

	if s.StartDate.After(time.Now()) {
		return errors.New("start_date cannot be in the future")
	}

	return nil
}
