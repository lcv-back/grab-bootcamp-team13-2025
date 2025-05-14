package app

import (
	"context"
	"fmt"
	"grab-bootcamp-be-team13-2025/internal/config"
	"grab-bootcamp-be-team13-2025/internal/domain/models"
	"grab-bootcamp-be-team13-2025/internal/domain/usecase"
	"grab-bootcamp-be-team13-2025/internal/handlers"
	"grab-bootcamp-be-team13-2025/internal/infrastructure/mysql"
	"grab-bootcamp-be-team13-2025/internal/infrastructure/email"
	"grab-bootcamp-be-team13-2025/internal/middleware"
	"io"
	applogger "grab-bootcamp-be-team13-2025/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
	mysqlGorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	jwt "grab-bootcamp-be-team13-2025/pkg/utils/jwt"
	"grab-bootcamp-be-team13-2025/pkg/utils/rabbitmq"
	"grab-bootcamp-be-team13-2025/pkg/utils/redis"
	customhttp "grab-bootcamp-be-team13-2025/pkg/utils/http"
	"github.com/bits-and-blooms/bloom/v3"
)

type App struct {
	Router         *gin.Engine
	Config         *config.Config
	DB             *gorm.DB
	redisClient    *redis.RedisClient
	rabbitmqClient *rabbitmq.RabbitMQClient
	version        string
	BloomFilter    *bloom.BloomFilter
	Logger         applogger.Logger
} // Thêm Logger vào struct App

// ProxyPredictHandler forwards requests to ngrok endpoint
func ProxyPredictHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := http.Post("https://cicada-logical-virtually.ngrok-free.app/predict", c.GetHeader("Content-Type"), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}
		c.Writer.WriteHeader(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}

func NewApp(configPath string, logger applogger.Logger) (*App, error) {
	// load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	//a.Logger.Info("JWT_SECRET:", cfg.JWTSecret)

	// Set Gin mode
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// connect to database with retries
	var db *gorm.DB
	retries := 10 // Tăng số lần thử lại lên 10
	for i := 0; i < retries; i++ {
		db, err = gorm.Open(mysqlGorm.Open(cfg.GetDSN()), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		})
		if err == nil {
			break
		}

		logger.Info("failed to connect to database", applogger.Field{Key: "host", Value: cfg.Database.Host}, applogger.Field{Key: "port", Value: cfg.Database.Port}, applogger.Field{Key: "retry", Value: i+1}, applogger.Field{Key: "max_retries", Value: retries}, applogger.Field{Key: "err", Value: err})
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database after %d retries: %w", retries, err)
	}

	// auto migrate tables
	if err = db.AutoMigrate(
		&models.User{},
		&models.Symptom{},
		&models.Disease{},
		&models.UserSymptom{},
		&models.PasswordResetToken{},
		&models.DiseaseSuggestion{},
		&models.DiseaseSymptom{},
		&models.Hospital{}, // migrate thêm bảng Hospital
	); err != nil {
		logger.Info("failed to auto-migrate tables", applogger.Field{Key: "err", Value: err})
		return nil, err
	}

	// Cập nhật toạ độ bệnh viện tự động 1 lần sau khi migrate
	if err := mysql.GeocodeAllHospitals(db); err != nil {
		logger.Info("Geocode hospitals failed", applogger.Field{Key: "err", Value: err})
	} else {
		logger.Info("Geocode hospitals completed!")
	}

	// initialize redis with retries
	var redisClient *redis.RedisClient
	for i := 0; i < retries; i++ {
		redisClient, err = redis.NewRedisClient(cfg.Redis.URL)
		if err == nil {
			break
		}

		logger.Info("failed to connect to redis", applogger.Field{Key: "retry", Value: i+1}, applogger.Field{Key: "max_retries", Value: retries}, applogger.Field{Key: "err", Value: err})
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis after %d retries: %v", retries, err)
	}

	// initialize rabbitmq with retries
	var rabbitmqClient *rabbitmq.RabbitMQClient
	for i := 0; i < retries; i++ {
		rabbitmqClient, err = rabbitmq.NewRabbitMQClient(cfg.RabbitMQ.URL)
		if err == nil {
			break
		}

		logger.Info("failed to connect to rabbitmq", applogger.Field{Key: "retry", Value: i+1}, applogger.Field{Key: "max_retries", Value: retries}, applogger.Field{Key: "err", Value: err})
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		if err = redisClient.Close(); err != nil {
			logger.Info("error closing redis connection", applogger.Field{Key: "err", Value: err})
		}
		return nil, fmt.Errorf("failed to initialize rabbitmq after %d retries: %v", retries, err)
	}

	// initialize repositories
	userRepo := mysql.NewMySQLUserRepository(db)
	symptomRepo := mysql.NewSymptomRepository(db) // Sửa tên hàm này
	diseaseRepo := mysql.NewMySQLDiseaseRepository(db)

	// initialize email service
	emailService := email.NewEmailService(os.Getenv("SENDGRID_API_KEY"))

	// Khởi tạo Bloom Filter
	// Ước lượng số lượng email: 10 triệu, xác suất false positive: 1%
	bloomFilter := bloom.NewWithEstimates(10_000_000, 0.01)

	// Tải tất cả email từ bảng users và thêm vào Bloom Filter
	var users []models.User
	if err := db.WithContext(context.Background()).Select("email").Find(&users).Error; err != nil {
		logger.Error("failed to load emails into Bloom Filter", applogger.Field{Key: "err", Value: err})
	}
	for _, user := range users {
		bloomFilter.Add([]byte(user.Email))
	}
	logger.Info("Initialized Bloom Filter", applogger.Field{Key: "emails", Value: len(users)})

	// initialize jwt utils
	jwtUtil := jwt.NewJWTUtil(cfg.JWTSecret)

	// initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, redisClient, emailService, bloomFilter)
	httpClient := customhttp.NewHTTPClient()
	symptomUsecase := usecase.NewSymptomUsecase(symptomRepo, diseaseRepo, redisClient, httpClient, cfg.MLService.URL)

	// initialize handlers
	authHandler := handlers.NewAuthHandler(authUsecase)
	symptomHandler := handlers.NewSymptomHandler(symptomUsecase, cfg)
	// Disease in-memory cache
	diseaseCache := handlers.NewDiseaseCache()
	{
		var diseases []models.Disease
		if err := db.WithContext(context.Background()).Find(&diseases).Error; err != nil {
			logger.Info("failed to load diseases into cache", applogger.Field{Key: "err", Value: err})
		} else {
			diseaseCache.Load(diseases)
		}
	}
	diseaseHandler := handlers.NewDiseaseHandler(diseaseRepo, diseaseCache)

	// initialize gin router
	// Add version info
	version := "1.0.0"

	// initialize gin router
	r := gin.New()

	// Add custom middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RateLimiter(redisClient))
	r.Use(middleware.TimeoutMiddleware(10 * time.Second))
	r.Use(middleware.RequestLogger())
	r.Use(middleware.PanicRecovery())

	// Health check endpoint with detailed status
	r.GET("/health", func(c *gin.Context) {
		// Check DB connection
		sqlDB, err := db.DB()
		dbStatus := "UP"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "DOWN"
		}

		// Check Redis connection
		redisStatus := "UP"
		if err := redisClient.Ping(context.Background()); err != nil {
			redisStatus = "DOWN"
		}

		// Check RabbitMQ connection
		rabbitmqStatus := "UP"
		if !rabbitmqClient.IsConnected() {
			rabbitmqStatus = "DOWN"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"version": version,
			"time":    time.Now().Format(time.RFC3339),
			"services": map[string]string{
				"database": dbStatus,
				"redis":    redisStatus,
				"rabbitmq": rabbitmqStatus,
			},
		})
	})

	// Hospital routes
	// Outbreak routes
	r.GET("/api/outbreaks", handlers.GetOutbreaksHandler)
	r.GET("/api/outbreaks/rss", handlers.GetOutbreaksFromRSSHandler)

	// Auth routes
	r.POST("/auth/signup", authHandler.Signup)
	r.POST("/auth/login", authHandler.Login)
	r.GET("/auth/me", middleware.JWTMiddleware(jwtUtil), authHandler.Me)
	r.POST("/auth/forgot-password", authHandler.ForgotPassword)
	r.POST("/auth/reset-password", authHandler.ResetPassword)
	r.PATCH("/api/update-info", middleware.JWTMiddleware(jwtUtil), authHandler.UpdateInfo)
	r.POST("/api/symptoms/upload", middleware.JWTMiddleware(jwtUtil), symptomHandler.Upload)
	r.GET("/api/symptoms/autocomplete", middleware.JWTMiddleware(jwtUtil), symptomHandler.AutocompleteSymptoms)
	r.POST("/api/symptoms/add", middleware.JWTMiddleware(jwtUtil), symptomHandler.AddSymptom)
	r.POST("/api/symptoms/predict", ProxyPredictHandler())

	// Route xóa ảnh MinIO
	r.DELETE("/api/images", middleware.JWTMiddleware(jwtUtil), symptomHandler.DeleteImage)

	// API tìm bệnh viện gần nhất bằng Google Places API (không dùng DB)
	r.POST("/api/hospitals/nearest", handlers.FindNearbyHospitalsExternal())

	// Disease API
	r.GET("/disease", diseaseHandler.GetDiseaseDescription)

	// Tái khởi tạo Bloom Filter định kỳ mỗi 24 giờ
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			// Tái khởi tạo Bloom Filter
			var users []models.User
			if err := db.WithContext(context.Background()).Select("email").Find(&users).Error; err != nil {
				logger.Info("failed to reload emails into Bloom Filter", applogger.Field{Key: "err", Value: err})
				continue
			}
			newBloomFilter := bloom.NewWithEstimates(10_000_000, 0.01)
			for _, user := range users {
				newBloomFilter.Add([]byte(user.Email))
			}
			bloomFilter = newBloomFilter
			logger.Info("Reinitialized Bloom Filter", applogger.Field{Key: "emails", Value: len(users)})
		}
	}()

	// Tạo queue RabbitMQ để đồng bộ Bloom Filter
	if err := rabbitmqClient.CreateQueue("bloom_filter_updates"); err != nil {
		logger.Error("failed to create bloom filter update queue", applogger.Field{Key: "err", Value: err})
	}

	// Consumer để cập nhật Bloom Filter từ RabbitMQ
	go func() {
		messages, err := rabbitmqClient.Consume("bloom_filter_updates")
		if err != nil {
			logger.Error("failed to consume bloom filter updates", applogger.Field{Key: "err", Value: err})
		}
		for msg := range messages {
			email := string(msg.Body)
			bloomFilter.Add([]byte(email))
			logger.Info("Updated Bloom Filter from RabbitMQ", applogger.Field{Key: "email", Value: email})
		}
	}()

	return &App{
		Router:         r,
		Config:         cfg,
		DB:             db,
		redisClient:    redisClient,
		rabbitmqClient: rabbitmqClient,
		version:        version,
		BloomFilter:    bloomFilter,
		Logger:         logger,
	}, nil
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.Config.Server.Port),
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		// Add TLS configuration if needed
		// TLSConfig: &tls.Config{
		//     MinVersion: tls.VersionTLS12,
		// },
	}

	// Start server in a goroutine
	go func() {
		a.Logger.Info("server started", applogger.Field{Key: "port", Value: srv.Addr})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("server listen error", applogger.Field{Key: "err", Value: err})
		}
	}()

	// Wait for interrupt signal
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Logger.Info("Shutting down server...")

	// Close connections
	if a.rabbitmqClient != nil {
		if err := a.rabbitmqClient.Close(); err != nil {
			a.Logger.Info("error closing rabbitmq connection", applogger.Field{Key: "err", Value: err})
		}
	}
	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			a.Logger.Info("error closing redis connection", applogger.Field{Key: "err", Value: err})
		}
	}
	if a.DB != nil {
		sqlDB, err := a.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				a.Logger.Info("error closing database connection", applogger.Field{Key: "err", Value: err})
			}
		}
	}

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		a.Logger.Error("server forced to shutdown", applogger.Field{Key: "err", Value: err})
		return err
	}

	// Close connections
	a.Shutdown()

	a.Logger.Info("server exited properly")
	return nil
}

func (a *App) Shutdown() {
	// Close Redis connection
	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			a.Logger.Info("error closing redis connection", applogger.Field{Key: "err", Value: err})
		}
	}

	// Close RabbitMQ connection
	if a.rabbitmqClient != nil {
		if err := a.rabbitmqClient.Close(); err != nil {
			a.Logger.Info("error closing rabbitmq connection", applogger.Field{Key: "err", Value: err})
		}
	}

	// Close DB connection
	if a.DB != nil {
		sqlDB, err := a.DB.DB()
		if err != nil {
			a.Logger.Info("error getting underlying sql.DB", applogger.Field{Key: "err", Value: err})
			return
		}
		if err := sqlDB.Close(); err != nil {
			a.Logger.Info("error closing database connection", applogger.Field{Key: "err", Value: err})
		}
	}
}
