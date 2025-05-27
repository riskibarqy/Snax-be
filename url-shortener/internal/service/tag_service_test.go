package service

import (
	"context"
	"testing"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTagRepository is a mock implementation of TagRepository
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Create(tag *domain.Tag) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *MockTagRepository) GetByID(id int64) (*domain.Tag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByName(name string) (*domain.Tag, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByURLID(urlID int64) ([]domain.Tag, error) {
	args := m.Called(urlID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tag), args.Error(1)
}

func (m *MockTagRepository) AddURLTag(urlID, tagID int64) error {
	args := m.Called(urlID, tagID)
	return args.Error(0)
}

func (m *MockTagRepository) RemoveURLTag(urlID, tagID int64) error {
	args := m.Called(urlID, tagID)
	return args.Error(0)
}

func (m *MockTagRepository) AddTagToURL(ctx context.Context, urlID int64, tag string) error {
	args := m.Called(ctx, urlID, tag)
	return args.Error(0)
}

func (m *MockTagRepository) RemoveTagFromURL(ctx context.Context, urlID int64, tag string) error {
	args := m.Called(ctx, urlID, tag)
	return args.Error(0)
}

func (m *MockTagRepository) GetURLTags(ctx context.Context, urlID int64) ([]domain.Tag, error) {
	args := m.Called(ctx, urlID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tag), args.Error(1)
}

func TestCreateTag(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tests := []struct {
		name      string
		tagName   string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:    "Success",
			tagName: "test-tag",
			mockSetup: func() {
				mockRepo.On("GetByName", "test-tag").Return(nil, assert.AnError)
				mockRepo.On("Create", mock.AnythingOfType("*domain.Tag")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Tag Already Exists",
			tagName: "existing-tag",
			mockSetup: func() {
				mockRepo.On("GetByName", "existing-tag").Return(&domain.Tag{
					ID:   1,
					Name: "existing-tag",
				}, nil)
			},
			wantErr: true,
		},
		{
			name:    "Repository Error",
			tagName: "test-tag",
			mockSetup: func() {
				mockRepo.On("GetByName", "test-tag").Return(nil, assert.AnError)
				mockRepo.On("Create", mock.AnythingOfType("*domain.Tag")).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			tag, err := service.CreateTag(tt.tagName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tag)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Equal(t, tt.tagName, tag.Name)
			}
		})
	}
}

func TestGetTag(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)

	tests := []struct {
		name      string
		tagID     int64
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "Success",
			tagID: 1,
			mockSetup: func() {
				mockRepo.On("GetByID", int64(1)).Return(&domain.Tag{
					ID:   1,
					Name: "test-tag",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Tag Not Found",
			tagID: 2,
			mockSetup: func() {
				mockRepo.On("GetByID", int64(2)).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			tag, err := service.GetTag(tt.tagID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tag)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tag)
				assert.Equal(t, tt.tagID, tag.ID)
			}
		})
	}
}

func TestGetURLTags(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)
	ctx := context.Background()

	tags := []domain.Tag{
		{ID: 1, Name: "tag1"},
		{ID: 2, Name: "tag2"},
	}

	tests := []struct {
		name      string
		urlID     int64
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "Success",
			urlID: 1,
			mockSetup: func() {
				mockRepo.On("GetURLTags", ctx, int64(1)).Return(tags, nil)
			},
			wantErr: false,
		},
		{
			name:  "Repository Error",
			urlID: 1,
			mockSetup: func() {
				mockRepo.On("GetURLTags", ctx, int64(1)).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			result, err := service.GetURLTags(ctx, tt.urlID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(tags), len(result))
				for i, tag := range result {
					assert.Equal(t, tags[i].ID, tag.ID)
					assert.Equal(t, tags[i].Name, tag.Name)
				}
			}
		})
	}
}

func TestAddTagToURL(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		urlID     int64
		tag       string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "Success",
			urlID: 1,
			tag:   "test-tag",
			mockSetup: func() {
				mockRepo.On("AddTagToURL", ctx, int64(1), "test-tag").Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Repository Error",
			urlID: 1,
			tag:   "test-tag",
			mockSetup: func() {
				mockRepo.On("AddTagToURL", ctx, int64(1), "test-tag").Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.AddTagToURL(ctx, tt.urlID, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveTagFromURL(t *testing.T) {
	mockRepo := new(MockTagRepository)
	service := NewTagService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		urlID     int64
		tag       string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "Success",
			urlID: 1,
			tag:   "test-tag",
			mockSetup: func() {
				mockRepo.On("RemoveTagFromURL", ctx, int64(1), "test-tag").Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Repository Error",
			urlID: 1,
			tag:   "test-tag",
			mockSetup: func() {
				mockRepo.On("RemoveTagFromURL", ctx, int64(1), "test-tag").Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.RemoveTagFromURL(ctx, tt.urlID, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
