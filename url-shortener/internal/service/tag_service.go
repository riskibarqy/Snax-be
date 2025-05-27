package service

import (
	"context"
	"strings"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type tagService struct {
	repo domain.TagRepository
}

// NewTagService creates a new tag service
func NewTagService(repo domain.TagRepository) domain.TagService {
	return &tagService{
		repo: repo,
	}
}

func (s *tagService) CreateTag(name string) (*domain.Tag, error) {
	// Normalize tag name
	name = strings.ToLower(strings.TrimSpace(name))

	// Check if tag already exists
	tag, err := s.repo.GetByName(name)
	if err == nil {
		return tag, nil
	}

	// Create new tag
	tag = &domain.Tag{
		Name: name,
	}
	err = s.repo.Create(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *tagService) GetTag(id int64) (*domain.Tag, error) {
	return s.repo.GetByID(id)
}

func (s *tagService) GetTagByName(name string) (*domain.Tag, error) {
	return s.repo.GetByName(strings.ToLower(strings.TrimSpace(name)))
}

func (s *tagService) GetURLTags(ctx context.Context, urlID int64) ([]domain.Tag, error) {
	return s.repo.GetURLTags(ctx, urlID)
}

func (s *tagService) AddTagToURL(ctx context.Context, urlID int64, tag string) error {
	return s.repo.AddTagToURL(ctx, urlID, tag)
}

func (s *tagService) RemoveTagFromURL(ctx context.Context, urlID int64, tag string) error {
	return s.repo.RemoveTagFromURL(ctx, urlID, tag)
}
