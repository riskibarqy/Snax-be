package service

import (
	"context"
	"testing"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockURLRepository is a mock implementation of URLRepository
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Create(url *domain.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) GetByShortCode(shortCode string) (*domain.URL, error) {
	args := m.Called(shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.URL), args.Error(1)
}

func (m *MockURLRepository) GetByUserID(userID string) ([]domain.URL, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.URL), args.Error(1)
}

func (m *MockURLRepository) Delete(id int64, userID string) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockURLRepository) IncrementClickCount(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateShortURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name        string
		originalURL string
		userID      string
		expiresAt   *time.Time
		mockSetup   func()
		wantErr     bool
	}{
		{
			name:        "Success",
			originalURL: "https://example.com",
			userID:      "user123",
			expiresAt:   nil,
			mockSetup: func() {
				mockRepo.On("Create", mock.AnythingOfType("*domain.URL")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Repository Error",
			originalURL: "https://example.com",
			userID:      "user123",
			expiresAt:   nil,
			mockSetup: func() {
				mockRepo.On("Create", mock.AnythingOfType("*domain.URL")).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			url, err := service.CreateShortURL(ctx, tt.originalURL, tt.userID, tt.expiresAt)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, url)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, url)
				assert.Equal(t, tt.originalURL, url.OriginalURL)
				assert.Equal(t, tt.userID, url.UserID)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	now := time.Now()
	expiredTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name      string
		shortCode string
		mockSetup func()
		wantErr   bool
		errType   interface{}
	}{
		{
			name:      "Success",
			shortCode: "abc123",
			mockSetup: func() {
				mockRepo.On("GetByShortCode", "abc123").Return(&domain.URL{
					ShortCode:   "abc123",
					OriginalURL: "https://example.com",
					ExpiresAt:   &futureTime,
					IsActive:    true,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "URL Not Found",
			shortCode: "notfound",
			mockSetup: func() {
				mockRepo.On("GetByShortCode", "notfound").Return(nil, &domain.ErrURLNotFound{ShortCode: "notfound"})
			},
			wantErr: true,
			errType: &domain.ErrURLNotFound{},
		},
		{
			name:      "URL Expired",
			shortCode: "expired",
			mockSetup: func() {
				mockRepo.On("GetByShortCode", "expired").Return(&domain.URL{
					ShortCode:   "expired",
					OriginalURL: "https://example.com",
					ExpiresAt:   &expiredTime,
					IsActive:    true,
				}, nil)
			},
			wantErr: true,
			errType: &domain.ErrURLExpired{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			url, err := service.GetURL(ctx, tt.shortCode)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, url)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, url)
				assert.Equal(t, tt.shortCode, url.ShortCode)
			}
		})
	}
}

func TestListUserURLs(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Success",
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("GetByUserID", "user123").Return([]domain.URL{
					{ID: 1, ShortCode: "abc123", UserID: "user123"},
					{ID: 2, ShortCode: "def456", UserID: "user123"},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Repository Error",
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("GetByUserID", "user123").Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			urls, err := service.ListUserURLs(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, urls)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, urls)
				for _, url := range urls {
					assert.Equal(t, tt.userID, url.UserID)
				}
			}
		})
	}
}

func TestDeleteURL(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		id        int64
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Success",
			id:     1,
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("Delete", int64(1), "user123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Repository Error",
			id:     1,
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("Delete", int64(1), "user123").Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.DeleteURL(ctx, tt.id, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRecordClick(t *testing.T) {
	mockRepo := new(MockURLRepository)
	service := NewURLService(mockRepo)

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
				mockRepo.On("IncrementClickCount", int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Repository Error",
			urlID: 1,
			mockSetup: func() {
				mockRepo.On("IncrementClickCount", int64(1)).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.RecordClick(tt.urlID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
