package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	AppEnv     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
}

func LoadConfig() *Config {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Println("Info: file .env not found, reading from system environment")
		envMap = make(map[string]string)
	}

	cfg := &Config{
		AppPort:    getEnvValue(envMap, "APP_PORT", "3000"),
		AppEnv:     getEnvValue(envMap, "APP_ENV", "development"),
		DBHost:     getEnvValue(envMap, "DB_HOST", "127.0.0.1"),
		DBUser:     getEnvValue(envMap, "DB_USER", "root"),
		DBPassword: getEnvValue(envMap, "DB_PASSWORD", ""),
		DBName:     getEnvValue(envMap, "DB_NAME", "ecommerce"),
		DBPort:     getEnvValue(envMap, "DB_PORT", "3306"),
		JWTSecret:  getEnvValue(envMap, "JWT_SECRET", ""),
	}

	return cfg
}

func (c *Config) ServerAddress() string {
	if !strings.Contains(c.AppPort, ":") {
		return ":" + c.AppPort
	}
	return c.AppPort
}

func getEnvValue(envMap map[string]string, key string, defaultValue string) string {
	if val, exists := envMap[key]; exists && val != "" {
		return val
	}
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// DSN
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// DB
func (c *Config) MaxIdleConns() int {
	return 10
}

func (c *Config) MaxOpenConns() int {
	return 25
}

func (c *Config) ConnMaxLifetime() time.Duration {
	return 5 * time.Minute
}

func (c *Config) ConnMaxIdleTime() time.Duration {
	return 1 * time.Minute
}
