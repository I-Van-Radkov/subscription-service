package adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewSubscriptionRepo(db *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *models.Subscription) (uuid.UUID, error) {
	query, args, err := r.builder.
		Insert("subscriptions").
		Columns("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		Values(sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.CreatedAt, sub.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	var id uuid.UUID
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert subscription: %w", err)
	}

	return id, nil
}

func (r *SubscriptionRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query, args, err := r.builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		From("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var sub models.Subscription
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price,
		&sub.UserID, &sub.StartDate, &sub.EndDate,
		&sub.CreatedAt, &sub.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscription: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub *models.Subscription) error {
	query, args, err := r.builder.
		Update("subscriptions").
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Set("updated_at", "NOW()").
		Where(squirrel.Eq{"id": sub.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	cmd, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.builder.
		Delete("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	cmd, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepo) List(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error) {
	qb := r.builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		From("subscriptions")

	if userID.String() != "" {
		qb = qb.Where(squirrel.Eq{"user_id": userID})
	}

	query, args, err := qb.OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	subs := make([]*models.Subscription, 0)
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(
			&s.ID, &s.ServiceName, &s.Price, &s.UserID,
			&s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		subs = append(subs, &s)
	}

	return subs, nil
}

func (r *SubscriptionRepo) SumForPeriod(ctx context.Context, userID uuid.UUID, serviceName string, start time.Time, end *time.Time) (int, error) {
	qb := r.builder.
		Select("COALESCE(SUM(price), 0) AS total").
		From("subscriptions").
		Where(squirrel.Or{
			squirrel.Expr("end_date IS NULL"),
			squirrel.GtOrEq{"end_date": start},
		})

	if end != nil {
		qb = qb.Where(squirrel.LtOrEq{"start_date": *end})
	}

	if userID.String() != "" {
		qb = qb.Where(squirrel.Eq{"user_id": userID})
	}

	if serviceName != "" {
		qb = qb.Where(squirrel.Eq{"service_name": serviceName})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build sum query: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get sum: %w", err)
	}

	return total, nil
}
