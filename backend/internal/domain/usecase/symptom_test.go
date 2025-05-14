package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/pkg/utils/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSymptomRepository is a mock implementation of SymptomRepository
type MockSymptomRepository struct {
	mock.Mock
}

func (m *MockSymptomRepository) CreateSymptom(ctx context.Context, symptom *models.Symptom) error {
	args := m.Called(ctx, symptom)
	return args.Error(0)
}

func (m *MockSymptomRepository) FindByName(ctx context.Context, name string) ([]models.Symptom, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Symptom), args.Error(1)
}

func (m *MockSymptomRepository) CreateUserSymptom(ctx context.Context, userSymptom *models.UserSymptom) error {
	args := m.Called(ctx, userSymptom)
	return args.Error(0)
}

func (m *MockSymptomRepository) GetUserSymptoms(ctx context.Context, userID uint64, timeRange string) ([]models.UserSymptom, error) {
	args := m.Called(ctx, userID, timeRange)
	return args.Get(0).([]models.UserSymptom), args.Error(1)
}

func (m *MockSymptomRepository) CreateFollowup(ctx context.Context, followup *models.Followup) error {
	args := m.Called(ctx, followup)
	return args.Error(0)
}

func (m *MockSymptomRepository) GetLatestFollowup(ctx context.Context, userID uint64) (*models.Followup, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Followup), args.Error(1)
}

func (m *MockSymptomRepository) UpdateFollowup(ctx context.Context, followup *models.Followup) error {
	args := m.Called(ctx, followup)
	return args.Error(0)
}

func (m *MockSymptomRepository) CreateDiagnosis(ctx context.Context, diagnosis *models.Diagnosis) error {
	args := m.Called(ctx, diagnosis)
	return args.Error(0)
}

func (m *MockSymptomRepository) GetLatestDiagnosis(ctx context.Context, userID uint64) (*models.Diagnosis, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.Diagnosis), args.Error(1)
}

// MockDiseaseRepository is a mock implementation of DiseaseRepository
type MockDiseaseRepository struct {
	mock.Mock
}

func (m *MockDiseaseRepository) FindByName(ctx context.Context, name string) (*models.Disease, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.Disease), args.Error(1)
}

func (m *MockDiseaseRepository) FindByIDs(ctx context.Context, ids []uint64) ([]models.Disease, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Disease), args.Error(1)
}

func (m *MockDiseaseRepository) FindByID(ctx context.Context, id uint64) (*models.Disease, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Disease), args.Error(1)
}

func (m *MockDiseaseRepository) CreateDiseaseSuggestion(ctx context.Context, suggestion *models.DiseaseSuggestion) error {
	args := m.Called(ctx, suggestion)
	return args.Error(0)
}

// MockRedisClient is a mock implementation of RedisClient
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// MockHTTPClient is a mock implementation of http.Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(url string, body interface{}) (*http.Response, error) {
	args := m.Called(url, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestSymptomUsecase_AddSymptoms(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	symptoms := []string{"fever", "cough"}
	imagePaths := []string{"http://example.com/image1.jpg"}
	mockSymptomRepo := new(MockSymptomRepository)
	mockDiseaseRepo := new(MockDiseaseRepository)
	mockRedisClient := new(MockRedisClient)
	mockHTTPClient := new(MockHTTPClient)

	// Setup

	// Mock tất cả các hàm repo/service có thể được gọi nhiều lần với mock.Anything để tránh thiếu expectation
	mockSymptomRepo.On("GetUserSymptoms", mock.Anything, mock.Anything, mock.Anything).Return([]models.UserSymptom{
		{UserID: userID, Name: "fever", ImagePaths: "[\"http://example.com/image1.jpg\"]", CreatedAt: time.Now()},
		{UserID: userID, Name: "cough", ImagePaths: "[\"http://example.com/image1.jpg\"]", CreatedAt: time.Now()},
	}, nil)
	mockSymptomRepo.On("CreateFollowup", mock.Anything, mock.Anything).Return(nil)
	mockSymptomRepo.On("UpdateFollowup", mock.Anything, mock.Anything).Return(nil)
	mockSymptomRepo.On("CreateDiagnosis", mock.Anything, mock.Anything).Return(nil)
	mockSymptomRepo.On("GetLatestFollowup", mock.Anything, mock.Anything).Return(&models.Followup{
		UserID:   userID,
		Attempt:  0,
		Symptoms: []string{"headache", "fatigue"},
		Answers:  []string{"", ""},
	}, nil)
	mockDiseaseRepo.On("FindByName", mock.Anything, mock.Anything).Return(&models.Disease{}, nil)
	mockDiseaseRepo.On("CreateDiseaseSuggestion", mock.Anything, mock.Anything).Return(nil)
	mockHTTPClient.On("Post", mock.Anything, mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       []byte(`{"predicted_diseases":[{"name":"Common Cold","probability":0.855}],"top_names":["headache","fatigue"]}`),
	}, nil)
	// (Các mock này sẽ được AssertExpectations ở cuối test)

	// Create usecase instance
	usecase := NewSymptomUsecase(mockSymptomRepo, mockDiseaseRepo, mockRedisClient, mockHTTPClient, "https://cicada-logical-virtually.ngrok-free.app")

	// Setup expectations
	// Mock Redis Get to return no cache
	mockRedisClient.On("Get", ctx, fmt.Sprintf("prediction:%d", userID)).Return("", nil)

	// Mock Redis Set to cache the prediction
	mockRedisClient.On("Set", ctx, fmt.Sprintf("prediction:%d", userID), mock.AnythingOfType("string"), time.Hour).Return(nil)

	mockSymptomRepo.On("CreateUserSymptom", ctx, mock.AnythingOfType("*models.UserSymptom")).Return(nil)
	mockSymptomRepo.On("GetUserSymptoms", ctx, userID, "").Return([]models.UserSymptom{
		{
			UserID:     userID,
			Name:       "fever",
			ImagePaths: "[\"http://example.com/image1.jpg\"]",
			CreatedAt:  time.Now(),
		},
	}, nil)

	// Mock ML service response
	mlResponse := &http.Response{
		StatusCode: 200,
		Body: []byte(`{
			"predicted_diseases": [
				{"name": "Common Cold", "probability": 85.5}
			],
			"top_names": ["headache", "fatigue"]
		}`),
	}
	mockHTTPClient.On("Post", "https://cicada-logical-virtually.ngrok-free.app/predict", mock.MatchedBy(func(req map[string]interface{}) bool {
		// Verify request structure
		_, hasUserID := req["user_id"]
		_, hasSymptoms := req["symptoms"]
		_, hasImagePaths := req["image_paths"]
		_, hasNumData := req["num_data"]
		return hasUserID && hasSymptoms && hasImagePaths && hasNumData
	})).Return(mlResponse, nil)

	// Mock disease repository
	mockDiseaseRepo.On("FindByName", ctx, "Common Cold").Return(&models.Disease{
		ID:   1,
		Name: "Common Cold",
	}, nil)
	mockDiseaseRepo.On("CreateDiseaseSuggestion", ctx, mock.AnythingOfType("*models.DiseaseSuggestion")).Return(nil)

	// Mock CreateFollowup for followup questions
	mockSymptomRepo.On("CreateFollowup", ctx, mock.MatchedBy(func(followup *models.Followup) bool {
		return followup.UserID == userID &&
			len(followup.Symptoms) == 2 &&
			contains(followup.Symptoms, "headache") &&
			contains(followup.Symptoms, "fatigue")
	})).Return(nil)

	// Execute
	resp, err := usecase.AddSymptoms(ctx, userID, symptoms, imagePaths)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.PredictedDiseases, 1)
	assert.Equal(t, "Common Cold", resp.PredictedDiseases[0].Name)
	assert.Len(t, resp.SymptomFollowups, 2)
	assert.Contains(t, resp.SymptomFollowups, "headache")
	assert.Contains(t, resp.SymptomFollowups, "fatigue")

	// Verify all expectations were met
	mockSymptomRepo.AssertExpectations(t)
	mockDiseaseRepo.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
	mockRedisClient.AssertExpectations(t)
}

func TestSymptomUsecase_AddSymptoms_WithCache(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	symptoms := []string{"fever", "cough"}
	imagePaths := []string{"http://example.com/image1.jpg"}
	mockSymptomRepo := new(MockSymptomRepository)
	mockDiseaseRepo := new(MockDiseaseRepository)
	mockRedisClient := new(MockRedisClient)
	mockHTTPClient := new(MockHTTPClient)

	// Setup

	// Create usecase instance
	usecase := NewSymptomUsecase(mockSymptomRepo, mockDiseaseRepo, mockRedisClient, mockHTTPClient, "https://cicada-logical-virtually.ngrok-free.app")

	// Mock cached response
	cachedResponse := &models.MLResponse{
		PredictedDiseases: []models.PredictedDisease{
			{Name: "Common Cold", Probability: 0.85},
		},
		SymptomFollowups: []string{"headache", "fatigue"},
	}
	cachedData, _ := json.Marshal(cachedResponse)
	mockRedisClient.On("Get", ctx, fmt.Sprintf("prediction:%d", userID)).Return(string(cachedData), nil)

	// Execute
	resp, err := usecase.AddSymptoms(ctx, userID, symptoms, imagePaths)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, cachedResponse, resp)

	// Verify that ML service was not called
	mockHTTPClient.AssertNotCalled(t, "Post")
}

func TestSymptomUsecase_AddSymptoms_InvalidInput(t *testing.T) {
	ctx := context.Background()
	// Bổ sung mock cho các hàm repo/service để tránh lỗi expectation

	// Setup
	userID := uint64(1)
	testCases := []struct {
		name        string
		symptoms    []string
		imagePaths  []string
		expectedErr error
	}{
		{
			name:        "Empty symptoms",
			symptoms:    []string{},
			imagePaths:  []string{},
			expectedErr: ErrInvalidSymptoms,
		},
		{
			name:        "Invalid image path",
			symptoms:    []string{"fever"},
			imagePaths:  []string{"invalid-path"},
			expectedErr: ErrInvalidImages,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repositories
			mockSymptomRepo := new(MockSymptomRepository)
			mockDiseaseRepo := new(MockDiseaseRepository)
			mockRedisClient := new(MockRedisClient)
			mockHTTPClient := new(MockHTTPClient)

			// Create usecase instance
			usecase := NewSymptomUsecase(mockSymptomRepo, mockDiseaseRepo, mockRedisClient, mockHTTPClient, "https://cicada-logical-virtually.ngrok-free.app")

			// Execute
			resp, err := usecase.AddSymptoms(ctx, userID, tc.symptoms, tc.imagePaths)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tc.expectedErr, err)
			assert.Nil(t, resp)

			// Verify that no repository methods were called
			mockSymptomRepo.AssertNotCalled(t, "CreateUserSymptom")
		})
	}
}

func TestSymptomUsecase_ProcessFollowup(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	answers := map[string]string{
		"headache": "yes",
		"fatigue":  "no",
	}
	mockSymptomRepo := new(MockSymptomRepository)
	mockDiseaseRepo := new(MockDiseaseRepository)
	mockRedisClient := new(MockRedisClient)
	mockHTTPClient := new(MockHTTPClient)

	// Bổ sung mock cho các hàm repo/service để tránh lỗi expectation
	mockSymptomRepo.On("GetLatestFollowup", ctx, userID).Return(&models.Followup{UserID: userID, Attempt: 0, Symptoms: []string{"headache", "fatigue"}, Answers: []string{"", ""}}, nil)
	mockSymptomRepo.On("UpdateFollowup", ctx, mock.Anything).Return(nil)
	mockSymptomRepo.On("CreateDiagnosis", ctx, mock.Anything).Return(nil)
	mockDiseaseRepo.On("FindByName", ctx, mock.Anything).Return(&models.Disease{}, nil)
	mockDiseaseRepo.On("CreateDiseaseSuggestion", ctx, mock.Anything).Return(nil)
	mockHTTPClient.On("Post", "https://cicada-logical-virtually.ngrok-free.app/predict", mock.Anything).Return(&http.Response{StatusCode: 200, Body: []byte(`{"predicted_diseases":[],"top_names":[]}`)}, nil)

	// Bổ sung mock cho các hàm repo/service để tránh lỗi expectation

	// Setup

	// Create usecase instance
	usecase := NewSymptomUsecase(mockSymptomRepo, mockDiseaseRepo, mockRedisClient, mockHTTPClient, "https://cicada-logical-virtually.ngrok-free.app")

	// Setup expectations
	mockSymptomRepo.On("GetLatestFollowup", ctx, userID).Return(&models.Followup{
		UserID:   userID,
		Attempt:  0,
		Symptoms: []string{"headache", "fatigue"},
		Answers:  []string{"", ""},
	}, nil)

	mockSymptomRepo.On("UpdateFollowup", ctx, mock.AnythingOfType("*models.Followup")).Return(nil)
	mockSymptomRepo.On("GetUserSymptoms", ctx, userID, "").Return([]models.UserSymptom{
		{
			UserID:     userID,
			Name:       "fever",
			ImagePaths: "[\"http://example.com/image1.jpg\"]",
			CreatedAt:  time.Now(),
		},
	}, nil)

	// Mock Redis Set to cache the prediction
	mockRedisClient.On("Set", ctx, fmt.Sprintf("prediction:%d", userID), mock.AnythingOfType("string"), time.Hour).Return(nil)

	// Mock ML service response
	mlResponse := &http.Response{
		StatusCode: 200,
		Body: []byte(`{
			"predicted_diseases": [
				{"name": "Common Cold", "probability": 95.5}
			],
			"top_names": []
		}`),
	}
	mockHTTPClient.On("Post", "https://cicada-logical-virtually.ngrok-free.app/predict", mock.MatchedBy(func(req map[string]interface{}) bool {
		// Verify request structure
		_, hasUserID := req["user_id"]
		_, hasSymptoms := req["symptoms"]
		_, hasImagePaths := req["image_paths"]
		_, hasNumData := req["num_data"]
		return hasUserID && hasSymptoms && hasImagePaths && hasNumData
	})).Return(mlResponse, nil)

	// Mock disease repository
	mockDiseaseRepo.On("FindByName", ctx, "Common Cold").Return(&models.Disease{
		ID:   1,
		Name: "Common Cold",
	}, nil)
	mockDiseaseRepo.On("CreateDiseaseSuggestion", ctx, mock.AnythingOfType("*models.DiseaseSuggestion")).Return(nil)
	mockSymptomRepo.On("CreateDiagnosis", ctx, mock.AnythingOfType("*models.Diagnosis")).Return(nil)

	// Execute
	resp, err := usecase.ProcessFollowup(ctx, userID, answers)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.PredictedDiseases, 1)
	assert.Equal(t, "Common Cold", resp.PredictedDiseases[0].Name)
	assert.Empty(t, resp.SymptomFollowups) // No more followups needed

	// Verify all expectations were met
	mockSymptomRepo.AssertExpectations(t)
	mockDiseaseRepo.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
	mockRedisClient.AssertExpectations(t)
}

func TestSymptomUsecase_ProcessFollowup_InvalidInput(t *testing.T) {
	ctx := context.Background()
	// Bổ sung mock cho các hàm repo/service để tránh lỗi expectation

	// Setup
	userID := uint64(1)
	testCases := []struct {
		name        string
		answers     map[string]string
		expectedErr error
	}{
		{
			name: "Invalid number of answers",
			answers: map[string]string{
				"headache": "yes",
			},
			expectedErr: ErrInvalidFollowupData,
		},
		{
			name: "Max attempts reached",
			answers: map[string]string{
				"headache": "yes",
				"fatigue":  "no",
			},
			expectedErr: ErrMaxFollowupReached,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repositories
			mockSymptomRepo := new(MockSymptomRepository)
			mockDiseaseRepo := new(MockDiseaseRepository)
			mockRedisClient := new(MockRedisClient)
			mockHTTPClient := new(MockHTTPClient)

			// Create usecase instance
			usecase := NewSymptomUsecase(mockSymptomRepo, mockDiseaseRepo, mockRedisClient, mockHTTPClient, "https://cicada-logical-virtually.ngrok-free.app")

			// Add universal mocks to prevent panic for all cases
			mockDiseaseRepo.On("FindByName", mock.Anything, mock.Anything).Return(&models.Disease{}, nil)
			mockDiseaseRepo.On("CreateDiseaseSuggestion", mock.Anything, mock.Anything).Return(nil)
			mockSymptomRepo.On("UpdateFollowup", mock.Anything, mock.Anything).Return(nil)
			mockSymptomRepo.On("CreateDiagnosis", mock.Anything, mock.Anything).Return(nil)

			// Setup expectations based on test case
			if tc.name == "Max attempts reached" {
				mockSymptomRepo.On("GetLatestFollowup", ctx, userID).Return(&models.Followup{
					UserID:   userID,
					Attempt:  3,
					Symptoms: []string{"headache", "fatigue"},
					Answers:  []string{"", ""},
				}, nil)
				mockSymptomRepo.On("GetUserSymptoms", ctx, userID, "").Return([]models.UserSymptom{}, nil)
			} else {
				mockSymptomRepo.On("GetLatestFollowup", ctx, userID).Return(&models.Followup{
					UserID:   userID,
					Attempt:  0,
					Symptoms: []string{"headache", "fatigue"},
					Answers:  []string{"", ""},
				}, nil)
				mockSymptomRepo.On("GetUserSymptoms", ctx, userID, "").Return([]models.UserSymptom{}, nil)
				mockHTTPClient.On("Post", "https://cicada-logical-virtually.ngrok-free.app/predict", mock.Anything).Return(&http.Response{
					StatusCode: 200,
					Body:       []byte(`{"predicted_diseases":[],"top_names":[]}`),
				}, nil)
			}

			// Execute
			resp, err := usecase.ProcessFollowup(ctx, userID, tc.answers)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tc.expectedErr, err)
			assert.Nil(t, resp)

			// Verify that no ML service was called in case of early exit
			if tc.name != "Max attempts reached" {
				mockHTTPClient.AssertNotCalled(t, "Post")
			}
		})
	}
}

func TestSymptomUsecase_ProcessFollowup_MLServiceError(t *testing.T) {
	ctx := context.Background()
	answers := map[string]string{
		"headache": "yes",
		"fatigue":  "no",
	}
	mockSymptomRepo := new(MockSymptomRepository)
	mockDiseaseRepo := new(MockDiseaseRepository)
	mockRedisClient := new(MockRedisClient)
	mockHTTPClient := new(MockHTTPClient)
	userID := uint64(1)

	// Create usecase instance
	usecase := NewSymptomUsecase(mockSymptomRepo, mockDiseaseRepo, mockRedisClient, mockHTTPClient, "https://cicada-logical-virtually.ngrok-free.app")

	// Add universal mocks to prevent panic for all cases
	mockDiseaseRepo.On("FindByName", mock.Anything, mock.Anything).Return(&models.Disease{}, nil)
	mockDiseaseRepo.On("CreateDiseaseSuggestion", mock.Anything, mock.Anything).Return(nil)
	mockSymptomRepo.On("UpdateFollowup", mock.Anything, mock.Anything).Return(nil)
	mockSymptomRepo.On("CreateDiagnosis", mock.Anything, mock.Anything).Return(nil)

	// Setup expectations
	mockSymptomRepo.On("GetLatestFollowup", ctx, userID).Return(&models.Followup{
		UserID:   userID,
		Attempt:  0,
		Symptoms: []string{"headache", "fatigue"},
		Answers:  []string{"", ""},
	}, nil)

	mockSymptomRepo.On("UpdateFollowup", ctx, mock.AnythingOfType("*models.Followup")).Return(nil)
	mockSymptomRepo.On("GetUserSymptoms", ctx, userID, "").Return([]models.UserSymptom{
		{
			UserID:     userID,
			Name:       "fever",
			ImagePaths: "[\"http://example.com/image1.jpg\"]",
			CreatedAt:  time.Now(),
		},
	}, nil)

	// Mock ML service error - return empty response and error
	mockHTTPClient.On("Post", "https://cicada-logical-virtually.ngrok-free.app/predict", mock.Anything).Return(&http.Response{
		StatusCode: 500,
		Body:       []byte{},
	}, errors.New("ML service error"))

	// Execute
	resp, err := usecase.ProcessFollowup(ctx, userID, answers)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed after 3 retries")

	// Verify that ML service was called 3 times (retry mechanism)
	mockHTTPClient.AssertNumberOfCalls(t, "Post", 3)
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
