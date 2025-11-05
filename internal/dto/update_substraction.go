package dto

type UpdateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"omitempty"`
	Price       int     `json:"price" validate:"omitempty,min=1"`
	StartDate   string  `json:"start_date" validate:"omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}
