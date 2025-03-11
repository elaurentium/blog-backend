package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/entities"
	"github.com/elaurentium/exilium-blog-backend/internal/domain/repositories"
	"github.com/google/uuid"
)

type SubService struct {
	subRepo repositories.SubRepository
	userRepo      repositories.UserRepository
}

func NewSubService(
	subRepo repositories.SubRepository,
	userRepo repositories.UserRepository,
) *SubService {
	return &SubService{
		subRepo: subRepo,
		userRepo:      userRepo,
	}
}

func (s *SubService) CreateSub(
	ctx context.Context,
	name string,
	description string,
	rules []string,
	creatorID uuid.UUID,
	isPrivate bool,
) (*entities.Sub, error) {
	// Verificar se o nome do sub é válido
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return nil, errors.New("sub name cannot be empty")
	}

	if len(name) < 3 || len(name) > 21 {
		return nil, errors.New("sub name must be between 3 and 21 characters")
	}

	// Verificar se o nome do sub já existe
	existingSub, err := s.subRepo.GetByName(ctx, name)
	if err == nil && existingSub != nil {
		return nil, errors.New("sub name already exists")
	}

	// Verificar se o criador existe
	_, err = s.userRepo.GetByID(ctx, creatorID)
	if err != nil {
		return nil, errors.New("creator not found")
	}

	now := time.Now()
	sub := &entities.Sub{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Rules:       rules,
		CreatorID:   creatorID,
		IsPrivate:   isPrivate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.subRepo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubService) GetSub(ctx context.Context, id uuid.UUID) (*entities.Sub, error) {
	return s.subRepo.GetByID(ctx, id)
}

func (s *SubService) GetSubByName(ctx context.Context, name string) (*entities.Sub, error) {
	return s.subRepo.GetByName(ctx, strings.ToLower(strings.TrimSpace(name)))
}

func (s *SubService) UpdateSub(
	ctx context.Context,
	id uuid.UUID,
	creatorID uuid.UUID,
	description string,
	rules []string,
	isPrivate bool,
	bannerURL string,
	iconURL string,
) (*entities.Sub, error) {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("sub not found")
	}

	// Verificar se o usuário é o criador do sub
	if sub.CreatorID != creatorID {
		return nil, errors.New("user not authorized to update this sub")
	}

	sub.Description = description
	sub.Rules = rules
	sub.IsPrivate = isPrivate
	sub.BannerURL = bannerURL
	sub.IconURL = iconURL
	sub.UpdatedAt = time.Now()

	err = s.subRepo.Update(ctx, sub)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubService) ListSubs(ctx context.Context, limit, offset int) ([]*entities.Sub, error) {
	return s.subRepo.List(ctx, limit, offset)
}

func (s *SubService) GetTrendingSub(ctx context.Context, limit int) ([]*entities.Sub, error) {
	return s.subRepo.GetTrending(ctx, limit)
}

func (s *SubService) DeleteSub(ctx context.Context, id uuid.UUID, creatorID uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("sub not found")
	}

	// Verificar se o usuário é o criador do sub
	if sub.CreatorID != creatorID {
		return errors.New("user not authorized to delete this sub")
	}

	return s.subRepo.Delete(ctx, id)
}