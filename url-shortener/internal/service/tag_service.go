package service

import (
	"context"
	"fmt"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type TagService struct {
	repo domain.TagRepository
}

// New creates a new tag service
func NewTagService(repo domain.TagRepository) domain.TagService {
	return &TagService{
		repo: repo,
	}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(name string) (*domain.Tag, error) {
	// Check if tag already exists
	existingTag, err := s.repo.GetByName(name)
	if err == nil && existingTag != nil {
		return nil, fmt.Errorf("tag already exists: %s", name)
	}

	tag := &domain.Tag{
		Name: name,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetTag retrieves a tag by ID
func (s *TagService) GetTag(id int64) (*domain.Tag, error) {
	return s.repo.GetByID(id)
}

// GetTagByName retrieves a tag by name
func (s *TagService) GetTagByName(name string) (*domain.Tag, error) {
	return s.repo.GetByName(name)
}

// GetURLTags retrieves all tags for a URL
func (s *TagService) GetURLTags(ctx context.Context, urlID int64) ([]domain.Tag, error) {
	return s.repo.GetURLTags(ctx, urlID)
}

// AddTagToURL adds a tag to a URL
func (s *TagService) AddTagToURL(ctx context.Context, urlID int64, tag string) error {
	return s.repo.AddTagToURL(ctx, urlID, tag)
}

// RemoveTagFromURL removes a tag from a URL
func (s *TagService) RemoveTagFromURL(ctx context.Context, urlID int64, tag string) error {
	return s.repo.RemoveTagFromURL(ctx, urlID, tag)
}
