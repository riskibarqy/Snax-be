package service

import (
	"context"
	"testing"
	"time"

	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnalyticsRepository is a mock implementation of AnalyticsRepository
type MockAnalyticsRepository struct {
	mock.Mock
}

func (m *MockAnalyticsRepository) Create(ctx context.Context, analytics *domain.Analytics) error {
	args := m.Called(ctx, analytics)
	return args.Error(0)
}

func (m *MockAnalyticsRepository) GetByURLID(ctx context.Context, urlID int64) ([]domain.Analytics, error) {
	args := m.Called(ctx, urlID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Analytics), args.Error(1)
}

func TestRecordVisit(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		urlID     int64
		visitorIP string
		userAgent string
		referer   string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Success",
			urlID:     1,
			visitorIP: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			referer:   "https://example.com",
			mockSetup: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Analytics")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Repository Error",
			urlID:     1,
			visitorIP: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			referer:   "https://example.com",
			mockSetup: func() {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Analytics")).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.RecordVisit(ctx, tt.urlID, tt.visitorIP, tt.userAgent, tt.referer)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetURLAnalytics(t *testing.T) {
	mockRepo := new(MockAnalyticsRepository)
	service := NewAnalyticsService(mockRepo)
	ctx := context.Background()

	now := time.Now()
	analytics := []domain.Analytics{
		{
			ID:          1,
			URLID:       1,
			VisitorIP:   "192.168.1.1",
			UserAgent:   "Mozilla/5.0",
			Referer:     "https://example.com",
			Timestamp:   now,
			CountryCode: "US",
			DeviceType:  "desktop",
		},
		{
			ID:          2,
			URLID:       1,
			VisitorIP:   "192.168.1.2",
			UserAgent:   "Mozilla/5.0",
			Referer:     "https://google.com",
			Timestamp:   now,
			CountryCode: "UK",
			DeviceType:  "mobile",
		},
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
				mockRepo.On("GetByURLID", ctx, int64(1)).Return(analytics, nil)
			},
			wantErr: false,
		},
		{
			name:  "Repository Error",
			urlID: 1,
			mockSetup: func() {
				mockRepo.On("GetByURLID", ctx, int64(1)).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			result, err := service.GetURLAnalytics(ctx, tt.urlID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(analytics), len(result))
				for i, a := range result {
					assert.Equal(t, analytics[i].ID, a.ID)
					assert.Equal(t, analytics[i].URLID, a.URLID)
					assert.Equal(t, analytics[i].VisitorIP, a.VisitorIP)
					assert.Equal(t, analytics[i].UserAgent, a.UserAgent)
					assert.Equal(t, analytics[i].Referer, a.Referer)
					assert.Equal(t, analytics[i].CountryCode, a.CountryCode)
					assert.Equal(t, analytics[i].DeviceType, a.DeviceType)
				}
			}
		})
	}
}
