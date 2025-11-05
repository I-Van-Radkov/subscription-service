package v1

import (
	"context"
	"net/http"

	"github.com/I-Van-Radkov/subscription-service/internal/dto"
	"github.com/gin-gonic/gin"
)

type SubscriptionUsecase interface {
	CreateSubscription(ctx context.Context, input dto.CreateSubstractionRequest) (dto.CreateSubstractionResponse, error)
	GetSubscription(ctx context.Context, id string) (dto.GetSubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, idString string, input dto.UpdateSubscriptionRequest) (dto.UpdateSubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id string) error
	GetSubscriptionsList(ctx context.Context, userId string) (dto.GetSubsListResponse, error)
	GetSubscriptionsSum(ctx context.Context, userIdStr, serviceName, start string, end *string) (dto.GetSubSumResponse, error)
}

type HandlerFacade struct {
	usecase SubscriptionUsecase
}

func NewHandlerFacade(usecase SubscriptionUsecase) *HandlerFacade {
	return &HandlerFacade{
		usecase: usecase,
	}
}

func (h *HandlerFacade) CreateSubscription(c *gin.Context) {
	var inputForm dto.CreateSubstractionRequest

	if err := c.ShouldBind(&inputForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outputForm, err := h.usecase.CreateSubscription(c.Request.Context(), inputForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, outputForm)
}

func (h *HandlerFacade) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	outputForm, err := h.usecase.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, outputForm)
}

func (h *HandlerFacade) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	var inputForm dto.UpdateSubscriptionRequest

	if err := c.ShouldBind(&inputForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outputForm, err := h.usecase.UpdateSubscription(c.Request.Context(), id, inputForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, outputForm)
}

func (h *HandlerFacade) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	err := h.usecase.DeleteSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "successful"})
}

func (h *HandlerFacade) GetSubscriptionsList(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	outputForm, err := h.usecase.GetSubscriptionsList(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, outputForm)
}

func (h *HandlerFacade) GetSubscriptionsSum(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	start := c.Query("start_date")
	end := c.Query("end_date")

	if start == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date is required"})
		return
	}
	outputForm, err := h.usecase.GetSubscriptionsSum(c.Request.Context(), userID, serviceName, start, &end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, outputForm)
}
