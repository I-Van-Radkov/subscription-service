package dto

type GetSubsListResponse struct {
	Total int                       `json:"total"`
	List  []GetSubscriptionResponse `json:"list"`
}
