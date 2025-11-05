package dto

type CreateSubstractionRequest struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Price       int     `json:"price" validate:"required,min=1"`
	UserID      string  `json:"user_id" validate:"required,uuid4"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

type CreateSubstractionResponse struct {
	ID string `json:"id" validate:"required,uuid4"`
}
