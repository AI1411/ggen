package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"g_gen/internal/domain/model"
	"g_gen/internal/usecase"
	mockdomain "g_gen/tests/mock/domain"
)

func setupPrefectureTest(t *testing.T) (*mockdomain.MockPrefectureRepository, usecase.PrefectureUseCase) {
	ctrl := gomock.NewController(t)
	mockRepo := mockdomain.NewMockPrefectureRepository(ctrl)
	useCase := usecase.NewPrefectureUseCase(mockRepo)
	return mockRepo, useCase
}

func TestPrefectureUseCase_ListPrefectures(t *testing.T) {
	// Setup
	mockRepo, useCase := setupPrefectureTest(t)
	ctx := context.Background()

	// Test cases
	tests := []struct {
		name          string
		mockSetup     func(mockRepo *mockdomain.MockPrefectureRepository)
		expectedError bool
		expectedLen   int
	}{
		{
			name: "Success",
			mockSetup: func(mockRepo *mockdomain.MockPrefectureRepository) {
				prefectures := []*model.Prefecture{
					{
						ID:   1,
						Name: "北海道",
					},
					{
						ID:   13,
						Name: "東京都",
					},
					{
						ID:   27,
						Name: "大阪府",
					},
				}
				mockRepo.EXPECT().FindAll(gomock.Any()).Return(prefectures, nil)
			},
			expectedError: false,
			expectedLen:   3,
		},
		{
			name: "Error",
			mockSetup: func(mockRepo *mockdomain.MockPrefectureRepository) {
				mockRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectedError: true,
			expectedLen:   0,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup(mockRepo)

			// Call the method
			prefectures, err := useCase.ListPrefectures(ctx)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, prefectures)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, prefectures)
				assert.Equal(t, tt.expectedLen, len(prefectures))
			}
		})
	}
}

func TestPrefectureUseCase_GetPrefectureByID(t *testing.T) {
	// Setup
	mockRepo, useCase := setupPrefectureTest(t)
	ctx := context.Background()

	// Test cases
	tests := []struct {
		name          string
		id            string
		mockSetup     func(mockRepo *mockdomain.MockPrefectureRepository)
		expectedError bool
	}{
		{
			name: "Success",
			id:   "13",
			mockSetup: func(mockRepo *mockdomain.MockPrefectureRepository) {
				prefecture := &model.Prefecture{
					ID:   13,
					Name: "東京都",
					Code: "13",
				}
				mockRepo.EXPECT().FindByCode(gomock.Any(), "13").Return(prefecture, nil)
			},
			expectedError: false,
		},
		{
			name: "Not Found",
			id:   "999",
			mockSetup: func(mockRepo *mockdomain.MockPrefectureRepository) {
				mockRepo.EXPECT().FindByCode(gomock.Any(), "999").Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			tt.mockSetup(mockRepo)

			// Call the method
			prefecture, err := useCase.GetPrefectureByID(ctx, tt.id)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, prefecture)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, prefecture)
				assert.Equal(t, tt.id, prefecture.Code)
			}
		})
	}
}
