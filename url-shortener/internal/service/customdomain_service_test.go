package service

import (
	"context"
	"testing"
	"time"

	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomDomainRepository is a mock implementation of CustomDomainRepository
type MockCustomDomainRepository struct {
	mock.Mock
}

func (m *MockCustomDomainRepository) Create(ctx context.Context, domain *internalDomain.CustomDomain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *MockCustomDomainRepository) GetByID(ctx context.Context, id int64) (*internalDomain.CustomDomain, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*internalDomain.CustomDomain), args.Error(1)
}

func (m *MockCustomDomainRepository) GetByDomain(ctx context.Context, domain string) (*internalDomain.CustomDomain, error) {
	args := m.Called(ctx, domain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*internalDomain.CustomDomain), args.Error(1)
}

func (m *MockCustomDomainRepository) GetByUserID(ctx context.Context, userID string) ([]internalDomain.CustomDomain, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]internalDomain.CustomDomain), args.Error(1)
}

func (m *MockCustomDomainRepository) Delete(ctx context.Context, id int64, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockCustomDomainRepository) VerifyDomain(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegisterDomain(t *testing.T) {
	mockRepo := new(MockCustomDomainRepository)
	service := NewCustomDomainService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		domain    string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Success",
			domain: "example.com",
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("GetByDomain", ctx, "example.com").Return(nil, &internalDomain.ErrDomainNotFound{Domain: "example.com"})
				mockRepo.On("Create", ctx, mock.MatchedBy(func(domain *internalDomain.CustomDomain) bool {
					return domain.Domain == "example.com" && domain.UserID == "user123" && !domain.Verified
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Domain Already Exists",
			domain: "existing.com",
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("GetByDomain", ctx, "existing.com").Return(&internalDomain.CustomDomain{
					ID:     1,
					Domain: "existing.com",
					UserID: "another-user",
				}, nil)
			},
			wantErr: true,
		},
		{
			name:   "Repository Error",
			domain: "example.com",
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("GetByDomain", ctx, "example.com").Return(nil, &internalDomain.ErrDomainNotFound{Domain: "example.com"})
				mockRepo.On("Create", ctx, mock.MatchedBy(func(domain *internalDomain.CustomDomain) bool {
					return domain.Domain == "example.com" && domain.UserID == "user123" && !domain.Verified
				})).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			domain, err := service.RegisterDomain(ctx, tt.domain, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, domain)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, domain)
				assert.Equal(t, tt.domain, domain.Domain)
				assert.Equal(t, tt.userID, domain.UserID)
				assert.False(t, domain.Verified)
			}
		})
	}
}

func TestGetUserDomains(t *testing.T) {
	mockRepo := new(MockCustomDomainRepository)
	service := NewCustomDomainService(mockRepo)
	ctx := context.Background()

	now := time.Now()
	domains := []internalDomain.CustomDomain{
		{
			ID:        1,
			Domain:    "example1.com",
			UserID:    "user123",
			Verified:  true,
			CreatedAt: now,
		},
		{
			ID:        2,
			Domain:    "example2.com",
			UserID:    "user123",
			Verified:  false,
			CreatedAt: now,
		},
	}

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
				mockRepo.On("GetByUserID", ctx, "user123").Return(domains, nil)
			},
			wantErr: false,
		},
		{
			name:   "Repository Error",
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("GetByUserID", ctx, "user123").Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			result, err := service.GetUserDomains(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(domains), len(result))
				for i, d := range result {
					assert.Equal(t, domains[i].ID, d.ID)
					assert.Equal(t, domains[i].Domain, d.Domain)
					assert.Equal(t, domains[i].UserID, d.UserID)
					assert.Equal(t, domains[i].Verified, d.Verified)
				}
			}
		})
	}
}

func TestDeleteDomain(t *testing.T) {
	mockRepo := new(MockCustomDomainRepository)
	service := NewCustomDomainService(mockRepo)
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
				mockRepo.On("Delete", ctx, int64(1), "user123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Repository Error",
			id:     1,
			userID: "user123",
			mockSetup: func() {
				mockRepo.On("Delete", ctx, int64(1), "user123").Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.DeleteDomain(ctx, tt.id, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyDomain(t *testing.T) {
	mockRepo := new(MockCustomDomainRepository)
	service := NewCustomDomainService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		id        int64
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success",
			id:   1,
			mockSetup: func() {
				mockRepo.On("VerifyDomain", ctx, int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Repository Error",
			id:   1,
			mockSetup: func() {
				mockRepo.On("VerifyDomain", ctx, int64(1)).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup()

			err := service.VerifyDomain(ctx, tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
