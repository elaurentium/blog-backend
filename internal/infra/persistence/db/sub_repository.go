package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/repositories"
)

type SubRepository struct {
	pool *pgxpool.Pool
}

func NewSubRepository(pool *pgxpool.Pool) repositories.SubRepository {
	return &SubRepository{pool: pool}
}

func (r *SubRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Sub, error) {
	query := `
		SELECT id, name, description, rules, creator_id, is_private, banner_url, icon_url, created_at, updated_at, deleted_at
		FROM subs
		WHERE id = $1 AND deleted_at IS NULL
	`

	var sub entities.Sub
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.Name, &sub.Description, &sub.Rules, &sub.CreatorID, &sub.IsPrivate, &sub.BannerURL, &sub.IconURL, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub by ID: %w", err)
	}

	return &sub, nil
}

func (r *SubRepository) GetByName(ctx context.Context, name string) (*entities.Sub, error) {
	query := `
		SELECT id, name, description, rules, creator_id, is_private, banner_url, icon_url, created_at, updated_at, deleted_at
		FROM subs
		WHERE name = $1 AND deleted_at IS NULL
	`

	var sub entities.Sub
	err := r.pool.QueryRow(ctx, query, name).Scan(
		&sub.ID, &sub.Name, &sub.Description, &sub.Rules, &sub.CreatorID, &sub.IsPrivate, &sub.BannerURL, &sub.IconURL, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub by name: %w", err)
	}

	return &sub, nil
}

func (r *SubRepository) Create(ctx context.Context, sub *entities.Sub) error {
	query := `
		INSERT INTO subs (id, name, description, rules, creator_id, is_private, banner_url, icon_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(ctx, query,
		sub.ID, sub.Name, sub.Description, sub.Rules, sub.CreatorID, sub.IsPrivate, sub.BannerURL, sub.IconURL, sub.CreatedAt, sub.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create sub: %w", err)
	}

	return nil
}

func (r *SubRepository) Update(ctx context.Context, sub *entities.Sub) error {
	query := `
		UPDATE subs
		SET name = $2, description = $3, rules = $4, is_private = $5, banner_url = $6, icon_url = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		sub.ID, sub.Name, sub.Description, sub.Rules, sub.IsPrivate, sub.BannerURL, sub.IconURL, sub.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update sub: %w", err)
	}

	return nil
}

func (r *SubRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE subs
		SET deleted_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub: %w", err)
	}

	return nil
}

func (r *SubRepository) List(ctx context.Context, limit, offset int) ([]*entities.Sub, error) {
	query := `
		SELECT id, name, description, rules, creator_id, is_private, banner_url, icon_url, created_at, updated_at
		FROM subs
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list subs: %w", err)
	}
	defer rows.Close()

	var subs []*entities.Sub
	for rows.Next() {
		var sub entities.Sub
		err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Description, &sub.Rules, &sub.CreatorID, &sub.IsPrivate, &sub.BannerURL, &sub.IconURL, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub: %w", err)
		}
		subs = append(subs, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over subs: %w", err)
	}

	return subs, nil
}

func (r *SubRepository) GetTrending(ctx context.Context, limit int) ([]*entities.Sub, error) {
	query := `
		SELECT s.id, s.name, s.description, s.rules, s.creator_id, s.is_private, s.banner_url, s.icon_url, s.created_at, s.updated_at
		FROM subs s
		LEFT JOIN posts p ON s.id = p.sub_id
		WHERE s.deleted_at IS NULL
		GROUP BY s.id
		ORDER BY COUNT(p.id) DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending subs: %w", err)
	}
	defer rows.Close()

	var subs []*entities.Sub
	for rows.Next() {
		var sub entities.Sub
		err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Description, &sub.Rules, &sub.CreatorID, &sub.IsPrivate, &sub.BannerURL, &sub.IconURL, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub: %w", err)
		}
		subs = append(subs, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over trending subs: %w", err)
	}

	return subs, nil
}