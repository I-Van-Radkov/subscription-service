package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/adapter"
	"github.com/I-Van-Radkov/subscription-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	srv *http.Server
	db  *pgxpool.Pool
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, db *pgxpool.Pool) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      nil,
	}

	return &Server{
		srv: srv,
		db:  db,
	}
}

func (s *Server) RegisterHandlers() error {
	subRepo := adapter.NewSubscriptionRepo(s.db)
	subUseCase := usecase.NewSubscriptionUsecase(subRepo)

	handler := NewHandlerFacade(subUseCase)

	router := gin.New()
	router.Use(LoggingMiddleware())

	api := router.Group("/api/v1")
	{
		api.POST("/subscriptions", handler.CreateSubscription)
		api.GET("/subscriptions/:id", handler.GetSubscription)
		api.PUT("/subscriptions/:id", handler.UpdateSubscription)
		api.DELETE("/subscriptions/:id", handler.DeleteSubscription)
		api.GET("/subscriptions", handler.GetSubscriptionsList)
		api.GET("/subscriptions/summary", handler.GetSubscriptionsSum)
	}

	s.srv.Handler = router

	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
