package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config структура, обозначающая структуру .env файла
type Config struct {
	Port              string // Порт, на котором запускается сервер
	MinioEndpoint     string // Адрес конечной точки Minio
	BucketName        string // Название конкретного бакета в Minio
	MinioRootUser     string // Имя пользователя для доступа к Minio
	MinioRootPassword string // Пароль для доступа к Minio
	MinioUseSSL       bool   // Переменная, отвечающая за
}

var AppConfig *Config

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Устанавливаем конфигурационные параметры
	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),

		MinioEndpoint:     getEnv("MINIO_ENDPOINT", "localhost:9001"),
		BucketName:        getEnv("MINIO_BUCKET_NAME", "mediaBucket"),
		MinioRootUser:     getEnv("MINIO_ROOT_USER", "blxxd"),
		MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", "minio_password"),
		MinioUseSSL:       getEnvAsBool("MINIO_USE_SSL", false),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
