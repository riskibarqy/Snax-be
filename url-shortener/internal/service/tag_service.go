package service

import (
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

func (s *tagService) GetURLTags(urlID int64) ([]domain.Tag, error) {
	return s.repo.GetByURLID(urlID)
}

func (s *tagService) AddTagToURL(urlID int64, tagName string) error {
	// Get or create tag
	tag, err := s.CreateTag(tagName)
	if err != nil {
		return err
	}

	// Add tag to URL
	return s.repo.AddURLTag(urlID, tag.ID)
}

func (s *tagService) RemoveTagFromURL(urlID int64, tagName string) error {
	// Get tag
	tag, err := s.GetTagByName(tagName)
	if err != nil {
		return err
	}

	// Remove tag from URL
	return s.repo.RemoveURLTag(urlID, tag.ID)
}
