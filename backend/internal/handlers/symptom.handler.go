// internal/handlers/symptom_handler.go
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/config"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/usecase"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	// Thêm các package cho decode/encode ảnh

	"image"
	"image/jpeg"


)

type SymptomHandler struct {
	symptomUsecase *usecase.SymptomUsecase
	s3Client       *s3.Client
	cfg            *config.Config
}

// DeleteImage handles DELETE /api/images with JSON body {"url": "<minio_url>"}
func (h *SymptomHandler) DeleteImage(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "url required in body"})
		return
	}
	// Parse filename or object key from URL
	urlParts := strings.Split(req.URL, "/")
	if len(urlParts) == 0 {
		c.JSON(400, gin.H{"error": "invalid url format"})
		return
	}
	filename := urlParts[len(urlParts)-1]
	if filename == "" {
		c.JSON(400, gin.H{"error": "could not extract filename from url"})
		return
	}
	_, err := h.s3Client.DeleteObject(c.Request.Context(), &s3.DeleteObjectInput{
		Bucket: aws.String(h.cfg.MinIO.Bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete image", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Image deleted successfully"})
}

func NewSymptomHandler(symptomUsecase *usecase.SymptomUsecase, cfg *config.Config) *SymptomHandler {
	// config minio using aws sdk
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: cfg.MinIO.Endpoint,
		}, nil
	})

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithEndpointResolverWithOptions(resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, "")),
		awsConfig.WithRegion("us-east-1"), // require but not use
	)
	if err != nil {
		log.Fatalf("failed to load MinIO config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // MinIO require path-style url
	})

	// check bucket minio exists
	_, err = s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(cfg.MinIO.Bucket),
	})
	if err != nil {
		log.Fatalf("failed to find bucket %s in MinIO: %v", cfg.MinIO.Bucket, err)
	}

	return &SymptomHandler{
		symptomUsecase: symptomUsecase,
		s3Client:       s3Client,
		cfg:            cfg,
	}
}

// AutocompleteSymptoms: suggest symptoms when user input text
func (h *SymptomHandler) AutocompleteSymptoms(c *gin.Context) {
	// get userID from context (saved by jwt middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	_, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	suggestions, err := h.symptomUsecase.AutocompleteSymptoms(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch symptom suggestions: %v", err)})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// AddSymptom handles the request to add a new symptom
func (h *SymptomHandler) AddSymptom(c *gin.Context) {
	// get userID from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse request body
	var req struct {
		Symptoms   []string          `json:"symptoms"`
		Answers    map[string]string `json:"answers"`
		ImagePaths []string          `json:"image_paths"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Convert answers (if present) to symptoms (all keys with value 'yes')
	symptoms := req.Symptoms
	if len(symptoms) == 0 && len(req.Answers) > 0 {
		for k, v := range req.Answers {
			if strings.ToLower(v) == "yes" {
				symptoms = append(symptoms, k)
			}
		}
	}

	// Validate image paths
	for _, path := range req.ImagePaths {
		if !strings.HasPrefix(path, "http") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image URL format"})
			return
		}
	}

	// Save symptoms to database
	for _, symptom := range symptoms {
		// Convert image paths to JSON string
		imagePathsJSON, err := json.Marshal(req.ImagePaths)
		if err != nil {
			log.Printf("Error marshaling image paths: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process image paths: %v", err)})
			return
		}

		userSymptom := &models.UserSymptom{
			UserID:     userID,
			Name:       symptom,
			ImagePaths: string(imagePathsJSON),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := h.symptomUsecase.CreateUserSymptom(c.Request.Context(), userSymptom); err != nil {
			log.Printf("Error saving symptom: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save symptom: %v", err)})
			return
		}
	}

	// Call ML service
	mlResp, err := h.symptomUsecase.AddSymptoms(c.Request.Context(), userID, symptoms, req.ImagePaths)
	if err != nil {
		log.Printf("Error calling ML service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get prediction: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":              "success",
		"message":             "Thêm triệu chứng thành công",
		"predicted_diseases":  mlResp.PredictedDiseases,
		"follow_up_questions": mlResp.SymptomFollowups,
	})
}

// CreateSymptom: save symptom, predict disease, and create follow-up questions
func (h *SymptomHandler) CreateSymptom(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ImagePath string `json:"image_path"`
	}

	// Bind data from form
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	// check name is required
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	// get userID from context (saved by jwt middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// process image upload (if exists)
	file, err := c.FormFile("image")
	var imagePath string
	if err == nil {
		log.Printf("Received image file: %s, size: %d", file.Filename, file.Size)
		// check file size
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image size exceeds 5MB"})
			return
		}

		// check file type
		ext := filepath.Ext(file.Filename)
		log.Printf("File extension: %s", ext)
		if !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".jpeg") && !strings.EqualFold(ext, ".png") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG images are allowed"})
			return
		}

		// create unique filename
		filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)

		// upload file to minio
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer fileReader.Close()

		_, err = h.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(h.cfg.MinIO.Bucket),
			Key:    aws.String(filename),
			Body:   fileReader,
			ACL:    "public-read",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload image to MinIO: %v", err)})
			return
		}

		// create public url for minio based on UseSSL
		scheme := "http"
		if h.cfg.MinIO.UseSSL {
			scheme = "https"
		}
		// remove protocol from endpoint to avoid duplicate
		endpoint := strings.TrimPrefix(h.cfg.MinIO.Endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
		imagePath = fmt.Sprintf("%s://%s/%s/%s", scheme, endpoint, h.cfg.MinIO.Bucket, filename)
	} else {
		log.Printf("No image file provided, using req.ImagePath: %s", req.ImagePath)
	}

	// Prepare image paths
	var imagePaths []string
	if imagePath != "" {
		imagePaths = append(imagePaths, imagePath)
	} else if req.ImagePath != "" {
		imagePaths = append(imagePaths, req.ImagePath)
	}

	// Convert image paths to JSON string
	imagePathsJSON, err := json.Marshal(imagePaths)
	if err != nil {
		log.Printf("Error marshaling image paths: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process image paths: %v", err)})
		return
	}

	// create UserSymptom object to save to database
	userSymptom := &models.UserSymptom{
		UserID:     userID,
		Name:       req.Name,
		ImagePaths: string(imagePathsJSON),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	log.Printf("Saving UserSymptom with ImagePaths: %s", userSymptom.ImagePaths)
	// save symptom to database
	if err = h.symptomUsecase.CreateUserSymptom(c.Request.Context(), userSymptom); err != nil {
		log.Printf("Failed to create user symptom: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create symptom: %v", err)})
		return
	}

	// predict disease and create follow-up questions
	top3Diseases, followUpQuestions, err := h.symptomUsecase.PredictDisease(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to predict diseases: %v", err)})
		return
	}

	// return success response
	c.JSON(http.StatusCreated, gin.H{
		"status":              "success",
		"message":             "Symptom created successfully",
		"predicted_diseases":  top3Diseases,
		"follow_up_questions": followUpQuestions,
	})
}

// Upload handles file upload for symptoms
func (h *SymptomHandler) Upload(c *gin.Context) {
	// get userID from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// limit file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10MB

	// parse multipart form with limit
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File too large. Maximum size is 10MB",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// check file type
	if !isValidFileType(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only JPG, JPEG, PNG, HEIC, WEBP files are allowed",
		})
		return
	}

	// Đọc file vào buffer để có thể decode nhiều lần
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	var uploadReader io.Reader
	var uploadFilename string

	switch ext {
	case ".heic":
		// Convert HEIC to JPEG
		// Decode HEIC to image.Image using standard image.Decode
		imgReader := bytes.NewReader(fileBytes)
		img, _, err := image.Decode(imgReader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode HEIC image: " + err.Error()})
			return
		}
		jpegBuf := new(bytes.Buffer)
		err = jpeg.Encode(jpegBuf, img, &jpeg.Options{Quality: 90})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode JPEG from HEIC"})
			return
		}
		uploadReader = jpegBuf
		uploadFilename = strings.TrimSuffix(header.Filename, ".heic") + ".jpg"

	default:
		// Các loại file khác giữ nguyên
		uploadReader = bytes.NewReader(fileBytes)
		uploadFilename = header.Filename
	}

	// create unique filename
	filename := fmt.Sprintf("symptoms/%d_%s_%s", userID, time.Now().Format("20060102150405"), uploadFilename)

	// upload file to minio
	_, err = h.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(h.cfg.MinIO.Bucket),
		Key:         aws.String(filename),
		Body:        uploadReader,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})

	if err != nil {
		log.Printf("Failed to upload file to MinIO: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file",
		})
		return
	}

	// create public url for file
	scheme := "http"
	if h.cfg.MinIO.UseSSL {
		scheme = "https"
	}
	publicEndpoint := "167.253.158.16:9000"
	fileURL := fmt.Sprintf("%s://%s/%s/%s", scheme, publicEndpoint, h.cfg.MinIO.Bucket, filename)
	log.Printf("Generated file URL: %s", fileURL)

	// Trả về URL để frontend có thể sử dụng khi tạo symptom
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fileURL,
	})
}

// isValidFileType check file type is valid or not
func isValidFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".heic": true,
	
	}
	return validTypes[ext]
}

// PredictSymptoms handles the request to predict diseases when the user clicks "Check"
func (h *SymptomHandler) PredictSymptoms(c *gin.Context) {
	// get userID from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse request body
	var req struct {
		UserID     string            `json:"user_id"`
		Symptoms   []string          `json:"symptoms"`
		Answers    map[string]string `json:"answers"`
		ImagePaths []string          `json:"image_paths"`
		NumData    int               `json:"num_data"`
	}

	// Debug: log raw request body
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	log.Printf("Raw request body: %s", string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON in PredictSymptoms: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Debug: log type of req.Answers
	log.Printf("Type of req.Answers: %T, value: %+v", req.Answers, req.Answers)

	// Nếu answers là nil thì gán map rỗng để tránh lỗi về sau
	if req.Answers == nil {
		req.Answers = map[string]string{}
	}

	// Nếu answers là nil thì gán map rỗng để tránh lỗi về sau
	if req.Answers == nil {
		req.Answers = map[string]string{}
	}

	// Convert answers (if present) to symptoms (all keys with value 'yes')
	symptoms := req.Symptoms
	if len(symptoms) == 0 && req.Answers != nil && len(req.Answers) > 0 {
		for k, v := range req.Answers {
			if strings.ToLower(v) == "yes" {
				symptoms = append(symptoms, k)
			}
		}
	}

	// Log the request for debugging
	log.Printf("Received prediction request for user %d: %+v", userID, req)

	// Validate image paths
	for _, path := range req.ImagePaths {
		if !strings.HasPrefix(path, "http") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image URL format"})
			return
		}
	}

	// Save symptoms to database
	for _, symptom := range symptoms {
		// Convert image paths to JSON string
		imagePathsJSON, err := json.Marshal(req.ImagePaths)
		if err != nil {
			log.Printf("Error marshaling image paths: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process image paths: %v", err)})
			return
		}

		userSymptom := &models.UserSymptom{
			UserID:     userID,
			Name:       symptom,
			ImagePaths: string(imagePathsJSON),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := h.symptomUsecase.CreateUserSymptom(c.Request.Context(), userSymptom); err != nil {
			log.Printf("Error saving symptom: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save symptom: %v", err)})
			return
		}
	}

	// call PredictedDisease to predict disease
	predictedDiseases, followUpQuestions, err := h.symptomUsecase.PredictDisease(c.Request.Context(), userID)
	if err != nil {
		// Log detailed error for debugging
		log.Printf("Error predicting diseases: %v", err)

		if strings.Contains(err.Error(), "invalid character") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "ML service đang gặp sự cố, vui lòng thử lại sau",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("không thể dự đoán bệnh: %v", err),
		})
		return
	}

	// Log successful prediction
	log.Printf("Successfully predicted diseases for user %d: %v", userID, predictedDiseases)

	c.JSON(http.StatusOK, gin.H{
		"status":              "success",
		"message":             "dự đoán bệnh thành công",
		"predicted_diseases":  predictedDiseases,
		"follow_up_questions": followUpQuestions,
	})
}

// AddSymptoms handles POST /api/symptoms/add
func (h *SymptomHandler) AddSymptoms(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	// Get symptoms
	symptoms := form.Value["symptoms"]
	if len(symptoms) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symptoms are required"})
		return
	}

	// Get images
	files := form.File["images"]
	imagePaths := make([]string, 0)
	for _, file := range files {
		// Validate file type
		if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only jpg and png images are allowed"})
			return
		}

		// Validate file size (5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image size must be less than 5MB"})
			return
		}

		// TODO: Upload to MinIO and get public URL
		// For now, just use the filename
		imagePaths = append(imagePaths, file.Filename)
	}

	// Process symptoms
	resp, err := h.symptomUsecase.AddSymptoms(c.Request.Context(), userID.(uint64), symptoms, imagePaths)
	if err != nil {
		switch err {
		case usecase.ErrInvalidSymptoms:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrInvalidImages:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrMLServerError:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ML service is unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Return response
	if len(resp.SymptomFollowups) == 0 {
		// Diagnosis complete
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"diagnosis": gin.H{
				"disease":    resp.PredictedDiseases[0].Name,
				"confidence": resp.PredictedDiseases[0].Probability,
				"details":    resp.Message,
			},
		})
	} else {
		// Need followup
		c.JSON(http.StatusOK, gin.H{
			"status":             "followup",
			"followup_questions": resp.SymptomFollowups,
			"message":            resp.Message,
		})
	}
}

// ProcessFollowup handles POST /api/symptoms/followup
func (h *SymptomHandler) ProcessFollowup(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse request body
	var req struct {
		Answers map[string]string `json:"answers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Process followup
	resp, err := h.symptomUsecase.ProcessFollowup(c.Request.Context(), userID.(uint64), req.Answers)
	if err != nil {
		switch err {
		case usecase.ErrMaxFollowupReached:
			c.JSON(http.StatusOK, gin.H{
				"status":  "incomplete",
				"message": "Unable to diagnose with high confidence. Please consult a doctor.",
			})
		case usecase.ErrMLServerError:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ML service is unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Return response
	if len(resp.SymptomFollowups) == 0 {
		// Diagnosis complete
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"diagnosis": gin.H{
				"disease":    resp.PredictedDiseases[0].Name,
				"confidence": resp.PredictedDiseases[0].Probability,
				"details":    resp.Message,
			},
		})
	} else {
		// Need followup
		c.JSON(http.StatusOK, gin.H{
			"status":             "followup",
			"followup_questions": resp.SymptomFollowups,
			"message":            resp.Message,
		})
	}
}
