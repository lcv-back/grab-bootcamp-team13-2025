package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/repository"
	"grab-bootcamp-be-team13-2025/pkg/utils/http"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidSymptoms      = errors.New("symptoms are required")
	ErrInvalidImages        = errors.New("invalid image format or size")
	ErrMaxFollowupReached   = errors.New("maximum followup attempts reached")
	ErrMLServerError        = errors.New("error from ML server")
	ErrInvalidFollowupData  = errors.New("invalid followup data")
	ErrMLServiceUnavailable = errors.New("ML service is unavailable")
	ErrMLServiceTimeout     = errors.New("ML service request timeout")
	ErrInvalidMLResponse    = errors.New("invalid response from ML service")
)

// MLRequest represents the request to ML service
type MLRequest struct {
	UserID     uint64   `json:"user_id"`
	Symptoms   []string `json:"symptoms"`
	ImagePaths []string `json:"image_paths"`
	NumData    int      `json:"num_data"`
}

// RedisClient defines the interface for Redis operations
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type SymptomUseCase interface {
	AddSymptoms(ctx context.Context, userID uint64, symptoms []string, imagePaths []string) (*models.MLResponse, error)
	ProcessFollowup(ctx context.Context, userID uint64, answers map[string]string) (*models.MLResponse, error)
	AutocompleteSymptoms(ctx context.Context, query string) ([]models.Symptom, error)
	CreateUserSymptom(ctx context.Context, userSymptom *models.UserSymptom) error
	PredictDisease(ctx context.Context, userID uint64) ([]models.PredictedDisease, []string, error)
}

type SymptomUsecase struct {
	symptomRepo  repository.SymptomRepository
	redisClient  RedisClient
	diseaseRepo  repository.DiseaseRepository
	mlClient     http.Client
	mlServiceURL string
	mu           sync.Mutex // Add mutex for concurrent access
}

func NewSymptomUsecase(symptomRepo repository.SymptomRepository, diseaseRepo repository.DiseaseRepository, redisClient RedisClient, mlClient http.Client, mlServiceURL string) *SymptomUsecase {
	return &SymptomUsecase{
		symptomRepo:  symptomRepo,
		diseaseRepo:  diseaseRepo,
		redisClient:  redisClient,
		mlClient:     mlClient,
		mlServiceURL: mlServiceURL,
	}
}

// validateSymptoms validates the input symptoms
func validateSymptoms(symptoms []string) error {
	if len(symptoms) == 0 {
		return ErrInvalidSymptoms
	}
	for _, s := range symptoms {
		if strings.TrimSpace(s) == "" {
			return ErrInvalidSymptoms
		}
	}
	return nil
}

// validateImages validates the input image paths
func validateImages(imagePaths []string) error {
	for _, path := range imagePaths {
		if !strings.HasPrefix(path, "http") {
			return ErrInvalidImages
		}
	}
	return nil
}

// getCachedPrediction gets cached prediction from Redis
func (u *SymptomUsecase) getCachedPrediction(ctx context.Context, userID uint64) (*models.MLResponse, error) {
	cacheKey := fmt.Sprintf("prediction:%d", userID)
	cached, err := u.redisClient.Get(ctx, cacheKey)
	if err == nil {
		var resp models.MLResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			log.Printf("Returning cached prediction for user %d", userID)
			return &resp, nil
		}
	}
	return nil, nil
}

// callMLServiceWithRetry calls ML service with retry mechanism
func (u *SymptomUsecase) callMLServiceWithRetry(ctx context.Context, req MLRequest) (*models.MLResponse, error) {
	var resp *models.MLResponse
	var err error
	for i := 0; i < 3; i++ {
		resp, err = u.callMLService(ctx, req)
		if err == nil {
			return resp, nil
		}
		log.Printf("Retry %d: ML service call failed: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return nil, fmt.Errorf("failed after 3 retries: %w", err)
}

// callMLService calls ML service and handles response
func (u *SymptomUsecase) callMLService(ctx context.Context, req MLRequest) (*models.MLResponse, error) {
	log.Printf("Calling ML service for user %d with %d symptoms", req.UserID, len(req.Symptoms))

	// Format request for ML service
	mlReq := map[string]interface{}{
		"user_id":     req.UserID,
		"symptoms":    req.Symptoms,
		"image_paths": req.ImagePaths,
		"num_data":    req.NumData,
	}

	// Log the request being sent
	log.Printf("Sending request to ML service: %+v", mlReq)

	// Call ML service
	//mlURL := strings.TrimRight(u.mlServiceURL, "/") + "/predict"
	mlURL := "https://cicada-logical-virtually.ngrok-free.app/predict"
	resp, err := u.mlClient.Post(mlURL, mlReq)
	if err != nil {
		log.Printf("Error calling ML service: %v", err)
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}

	// Log the response received
	log.Printf("Received response from ML service with status code: %d", resp.StatusCode)

	// Check status code
	if resp.StatusCode != 200 {
		log.Printf("ML service returned error status: %d", resp.StatusCode)
		return nil, fmt.Errorf("ML service returned error status: %d", resp.StatusCode)
	}

	// Parse response
	var rawResp map[string]interface{}
	if err = json.Unmarshal(resp.Body, &rawResp); err != nil {
		log.Printf("Failed to decode ML response. Body: %s, Error: %v", string(resp.Body), err)
		return nil, fmt.Errorf("failed to decode ML response: %w", err)
	}

	// Log the parsed response
	log.Printf("Successfully parsed ML response: %+v", rawResp)

	// Convert to MLResponse
	mlResp := &models.MLResponse{
		PredictedDiseases: make([]models.PredictedDisease, 0),
		SymptomFollowups:  make([]string, 0),
	}

	// Parse predicted diseases
	if diseases, ok := rawResp["predicted_diseases"].([]interface{}); ok {
		for _, d := range diseases {
			if disease, ok := d.(map[string]interface{}); ok {
				if name, ok := disease["name"].(string); ok {
					if prob, ok := disease["probability"].(float64); ok {
						mlResp.PredictedDiseases = append(mlResp.PredictedDiseases, models.PredictedDisease{
							Name:        name,
							Probability: prob / 100, // Convert percentage to decimal
						})
					}
				}
			}
		}
	}

	// Parse followup questions
	if topNames, ok := rawResp["top_names"].([]interface{}); ok {
		for _, name := range topNames {
			if symptom, ok := name.(string); ok {
				mlResp.SymptomFollowups = append(mlResp.SymptomFollowups, symptom)
			}
		}
	}

	// Log the final response
	log.Printf("Final ML response: %+v", mlResp)

	// Cache the response
	if len(mlResp.PredictedDiseases) > 0 {
		cacheKey := fmt.Sprintf("prediction:%d", req.UserID)
		cachedData, err := json.Marshal(mlResp)
		if err == nil {
			if err := u.redisClient.Set(ctx, cacheKey, string(cachedData), time.Hour); err != nil {
				log.Printf("Failed to cache prediction for user %d: %v", req.UserID, err)
			}
		}
	}

	return mlResp, nil
}

// filterExistingSymptoms filters out symptoms that user already reported
func (u *SymptomUsecase) filterExistingSymptoms(ctx context.Context, userID uint64, followups []string) ([]string, error) {
	existing, err := u.symptomRepo.GetUserSymptoms(ctx, userID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get existing symptoms: %w", err)
	}

	existingMap := make(map[string]bool)
	for _, s := range existing {
		existingMap[s.Name] = true
	}

	filtered := make([]string, 0)
	for _, f := range followups {
		if !existingMap[f] {
			filtered = append(filtered, f)
		}
	}
	return filtered, nil
}

// saveDiseaseSuggestions saves disease suggestions to the database
func (u *SymptomUsecase) saveDiseaseSuggestions(ctx context.Context, userID uint64, diseases []models.PredictedDisease) error {
	for _, disease := range diseases {
		// Get disease ID
		diseaseModel, err := u.diseaseRepo.FindByName(ctx, disease.Name)
		if err != nil {
			return fmt.Errorf("failed to find disease: %w", err)
		}

		// Create disease suggestion
		suggestion := &models.DiseaseSuggestion{
			UserID:    userID,
			DiseaseID: diseaseModel.ID,
			CreatedAt: time.Now(),
		}

		if err := u.diseaseRepo.CreateDiseaseSuggestion(ctx, suggestion); err != nil {
			return fmt.Errorf("failed to create disease suggestion: %w", err)
		}
	}
	return nil
}

// AddSymptoms adds new symptoms for a user
func (u *SymptomUsecase) AddSymptoms(ctx context.Context, userID uint64, symptoms []string, imagePaths []string) (*models.MLResponse, error) {
	// Validate input
	if err := validateSymptoms(symptoms); err != nil {
		return nil, err
	}
	if err := validateImages(imagePaths); err != nil {
		return nil, err
	}

	// Check cache first
	if cached, err := u.getCachedPrediction(ctx, userID); err == nil && cached != nil {
		log.Printf("Returning cached prediction for user %d", userID)
		return cached, nil
	}

	// Save symptoms to database
	for _, symptom := range symptoms {
		// Convert image paths to JSON string
		imagePathsJSON, err := json.Marshal(imagePaths)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal image paths: %w", err)
		}

		userSymptom := &models.UserSymptom{
			UserID:     userID,
			Name:       symptom,
			ImagePaths: string(imagePathsJSON),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := u.symptomRepo.CreateUserSymptom(ctx, userSymptom); err != nil {
			return nil, fmt.Errorf("failed to save symptom: %w", err)
		}
	}

	// Call ML service
	mlReq := MLRequest{
		UserID:     userID,
		Symptoms:   symptoms,
		ImagePaths: imagePaths,
		NumData:    5,
	}

	mlResp, err := u.callMLServiceWithRetry(ctx, mlReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	// Cache the response
	if len(mlResp.PredictedDiseases) > 0 {
		cacheKey := fmt.Sprintf("prediction:%d", userID)
		cachedData, err := json.Marshal(mlResp)
		if err == nil {
			if err := u.redisClient.Set(ctx, cacheKey, string(cachedData), time.Hour); err != nil {
				log.Printf("Failed to cache prediction: %v", err)
			}
		}
	}

	return mlResp, nil
}

// ProcessFollowup processes followup answers
func (u *SymptomUsecase) ProcessFollowup(ctx context.Context, userID uint64, answers map[string]string) (*models.MLResponse, error) {
	// Get user symptoms
	symptoms, err := u.symptomRepo.GetUserSymptoms(ctx, userID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user symptoms: %w", err)
	}

	// Collect all symptoms and image paths
	var allSymptoms []string
	var allImagePaths []string
	for _, s := range symptoms {
		allSymptoms = append(allSymptoms, s.Name)

		// Parse image paths from JSON string
		var imagePaths []string
		paths := s.ImagePaths
		if paths == "" {
			paths = "[]"
		}
		if err := json.Unmarshal([]byte(paths), &imagePaths); err != nil {
			log.Printf("Failed to unmarshal image paths: %v", err)
			return nil, fmt.Errorf("failed to unmarshal image paths: %w", err)
		}
		allImagePaths = append(allImagePaths, imagePaths...)
	}

	// Call ML service
	mlReq := MLRequest{
		UserID:     userID,
		Symptoms:   allSymptoms,
		ImagePaths: allImagePaths,
		NumData:    5,
	}

	mlResp, err := u.callMLServiceWithRetry(ctx, mlReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	return mlResp, nil
}

// AutocompleteSymptoms returns symptom suggestions based on query
func (u *SymptomUsecase) AutocompleteSymptoms(ctx context.Context, query string) ([]models.Symptom, error) {
	return u.symptomRepo.FindByName(ctx, query)
}

// CreateUserSymptom creates a new user symptom
func (u *SymptomUsecase) CreateUserSymptom(ctx context.Context, userSymptom *models.UserSymptom) error {
	// Find symptom ID from name
	symptoms, err := u.symptomRepo.FindByName(ctx, userSymptom.Name)
	if err != nil {
		return fmt.Errorf("failed to find symptom: %v", err)
	}

	var symptom models.Symptom
	if len(symptoms) == 0 {
		// If symptom not found, create new symptom
		symptom = models.Symptom{
			Name: userSymptom.Name,
		}
		if err := u.symptomRepo.CreateSymptom(ctx, &symptom); err != nil {
			return fmt.Errorf("failed to create symptom: %v", err)
		}
	} else {
		symptom = symptoms[0]
	}

	// Set symptom ID
	userSymptom.SymptomID = symptom.ID

	// Create user symptom
	if err := u.symptomRepo.CreateUserSymptom(ctx, userSymptom); err != nil {
		return fmt.Errorf("failed to create user symptom: %v", err)
	}

	return nil
}

// PredictDisease predicts diseases based on user symptoms
func (u *SymptomUsecase) PredictDisease(ctx context.Context, userID uint64) ([]models.PredictedDisease, []string, error) {
	// Get user symptoms
	symptoms, err := u.symptomRepo.GetUserSymptoms(ctx, userID, "")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user symptoms: %w", err)
	}

	// Collect all symptoms and image paths
	var allSymptoms []string
	var allImagePaths []string
	for _, s := range symptoms {
		allSymptoms = append(allSymptoms, s.Name)

		// Parse image paths from JSON string
		var imagePaths []string
		paths := s.ImagePaths
		if paths == "" {
			paths = "[]"
		}
		if err := json.Unmarshal([]byte(paths), &imagePaths); err != nil {
			log.Printf("Failed to unmarshal image paths: %v", err)
			continue
		}
		allImagePaths = append(allImagePaths, imagePaths...)
	}

	// Call ML service
	mlReq := MLRequest{
		UserID:     userID,
		Symptoms:   allSymptoms,
		ImagePaths: allImagePaths,
		NumData:    5,
	}

	mlResp, err := u.callMLServiceWithRetry(ctx, mlReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get prediction: %w", err)
	}

	return mlResp.PredictedDiseases, mlResp.SymptomFollowups, nil
}
