package db

import (
	"context"
	"fmt"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, username, email, hashed_password, salt, birthday, full_name, bio, avatar_url, role, is_active, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.HashedPassword, user.Salt, user.Birthday, user.FullName,
		user.Bio, user.AvatarURL, user.Role, user.IsActive, user.EmailVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, username, email, hashed_password, salt, full_name, bio, avatar_url, role, is_active, email_verified, last_login, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Salt, &user.FullName,
		&user.Bio, &user.AvatarURL, &user.Role, &user.IsActive, &user.EmailVerified, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err // Handle sql.ErrNoRows as needed
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, username, email, hashed_password, salt, full_name, bio, avatar_url, role, is_active, email_verified, last_login, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Salt, &user.FullName,
		&user.Bio, &user.AvatarURL, &user.Role, &user.IsActive, &user.EmailVerified, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err // Handle sql.ErrNoRows as needed
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users SET username = $2, email = $3, hashed_password = $4, salt = $5, full_name = $6, bio = $7,
			avatar_url = $8, role = $9, is_active = $10, email_verified = $11, last_login = $12, updated_at = $13
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.HashedPassword, user.Salt, user.FullName,
		user.Bio, user.AvatarURL, user.Role, user.IsActive, user.EmailVerified, user.LastLogin, user.UpdatedAt)
	return err
}

func (r *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)", email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return exists, nil
}

func (r *userRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL)", username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return exists, nil
}