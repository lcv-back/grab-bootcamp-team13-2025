// internal/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Database struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Name      string `yaml:"name"`
		Charset   string `yaml:"charset"`
		ParseTime bool   `yaml:"parse_time"`
		Loc       string `yaml:"loc"`
	} `yaml:"database"`
	Email struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"email"`
	JWTSecret string
	Redis     struct {
		URL string `yaml:"url"`
	} `yaml:"redis"`
	RabbitMQ struct {
		URL string `yaml:"url"`
	} `yaml:"rabbitmq"`
	MinIO struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucket"`
		UseSSL    bool   `yaml:"use_ssl"`
	} `yaml:"minio"`
	MLService struct {
		URL string `yaml:"url"`
	} `yaml:"ml_service"`
	ResetPasswordURL string
	RedisURL         string
	RabbitMQURL      string
	SendGridAPIKey string
}

func LoadConfig(configPath string) (*Config, error) {
	// Load config from yaml file
	config := &Config{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	config.SendGridAPIKey = os.Getenv("SENDGRID_API_KEY")

	// Load JWT_SECRET from .env
	config.JWTSecret = os.Getenv("JWT_SECRET")
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET not found in .env")
	}

	// Load database config from .env if not set in config file
	if config.Database.Host == "" {
		config.Database.Host = os.Getenv("DB_HOST")
		if config.Database.Host == "" {
			log.Fatal("DB_HOST not found in .env")
		}
	}
	if config.Database.Port == 0 {
		portStr := os.Getenv("DB_PORT")
		port, err := strconv.Atoi(portStr)
		if err != nil || port == 0 {
			log.Fatal("DB_PORT invalid in .env")
		}
		config.Database.Port = port
	}
	if config.Database.Name == "" {
		config.Database.Name = os.Getenv("DB_NAME")
		if config.Database.Name == "" {
			log.Fatal("DB_NAME not found in .env")
		}
	}
	if config.Database.User == "" {
		config.Database.User = os.Getenv("DB_USER")
		if config.Database.User == "" {
			log.Fatal("DB_USER not found in .env")
		}
	}
	if config.Database.Password == "" {
		config.Database.Password = os.Getenv("DB_PASSWORD")
		if config.Database.Password == "" {
			log.Fatal("DB_PASSWORD not found in .env")
		}
	}

	// Load email config from environment variables
	config.Email.Host = os.Getenv("EMAIL_HOST")
	if config.Email.Host == "" {
		log.Fatal("EMAIL_HOST not found in .env")
	}

	portStr := os.Getenv("EMAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		log.Fatal("EMAIL_PORT invalid in .env")
	}
	config.Email.Port = port

	config.Email.Username = os.Getenv("EMAIL_USERNAME")
	if config.Email.Username == "" {
		log.Fatal("EMAIL_USERNAME not found in .env")
	}

	config.Email.Password = os.Getenv("EMAIL_PASSWORD")
	if config.Email.Password == "" {
		log.Fatal("EMAIL_PASSWORD not found in .env")
	}

	// Load MinIO config from environment variables if not set in config file
	// Ưu tiên giá trị từ config.yaml, chỉ đọc từ .env nếu giá trị trong config.yaml rỗng
	if config.MinIO.Endpoint == "" {
		config.MinIO.Endpoint = os.Getenv("MINIO_ENDPOINT")
	}
	if config.MinIO.AccessKey == "" {
		config.MinIO.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	}
	if config.MinIO.SecretKey == "" {
		config.MinIO.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	}
	if config.MinIO.Bucket == "" {
		config.MinIO.Bucket = os.Getenv("MINIO_BUCKET")
	}

	// UseSSL mặc định là false nếu không được đặt trong YAML hoặc .env
	if !config.MinIO.UseSSL {
		useSSLStr := os.Getenv("MINIO_USE_SSL")
		if useSSLStr != "" {
			useSSL, err := strconv.ParseBool(useSSLStr)
			if err != nil {
				log.Fatal("MINIO_USE_SSL invalid in .env (must be true/false)")
			}
			config.MinIO.UseSSL = useSSL
		}
	}

	// Kiểm tra các thông số MinIO bắt buộc
	if config.MinIO.Endpoint == "" || config.MinIO.AccessKey == "" || config.MinIO.SecretKey == "" || config.MinIO.Bucket == "" {
		log.Fatal("MinIO configuration incomplete: endpoint, access_key, secret_key, and bucket are required")
	}

	config.ResetPasswordURL = getEnv("RESET_PASSWORD_URL", "https://isymptom.vercel.app/reset-password")
	config.RedisURL = getEnv("REDIS_URL", "167.253.158.16:6379")
	config.RabbitMQURL = getEnv("RABBITMQ_URL", "amqp://admin:admin@167.253.158.16:5672/")

	return config, nil
}

// GetDSN return DSN to connect to MySQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.Charset,
		c.Database.ParseTime,
		c.Database.Loc,
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
