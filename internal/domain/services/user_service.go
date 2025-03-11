package services

import (
	"context"
	"errors"
	"time"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/repositories"
	"github.com/elaurentium/exilium-blog-backend/internal/infra/auth"
	"github.com/google/uuid"
)
type UserService struct {
	userRepo repositories.UserRepository
	auth     auth.AuthService
}

func NewUserService(userRepo repositories.UserRepository, auth auth.AuthService) *UserService {
	return &UserService{
		userRepo: userRepo,
		auth:     auth,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password, fullName string, birthday string) (*entities.User, error) {
	emailExists, err := s.userRepo.CheckEmailExists(ctx, email)

	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, salt, err := s.auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &entities.User{
		ID:             uuid.New(),
		Username:       username,
		Email:          email,
		Birthday:       birthday,
		HashedPassword: hashedPassword,
		Salt:           salt,
		FullName:       fullName,
		Role:           "user",
		IsActive:       true,
		EmailVerified:  false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !s.auth.VerifyPassword(password, user.HashedPassword, user.Salt) {
		return "", "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", "", errors.New("account is deactivated")
	}

	now := time.Now()
	user.LastLogin = &now
	user.UpdatedAt = now
	
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, fullName, bio, avatarURL string) (*entities.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = fullName
	user.Bio = bio
	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !s.auth.VerifyPassword(currentPassword, user.HashedPassword, user.Salt) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, salt, err := s.auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.HashedPassword = hashedPassword
	user.Salt = salt
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}