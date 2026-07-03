package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string        `env:"APP_ENV" envDefault:"development"`
	HTTPAddr        string        `env:"HTTP_ADDR" envDefault:":8080"`
	DatabaseURL     string        `env:"DATABASE_URL,required"`
	RedisURL        string        `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`
	JWTAccessSecret string        `env:"JWT_ACCESS_SECRET,required"`
	JWTRefreshSecret string       `env:"JWT_REFRESH_SECRET,required"`
	JWTAccessTTL    time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	JWTRefreshTTL   time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h"`
	BcryptCost      int           `env:"BCRYPT_COST" envDefault:"12"`
	AllowedOrigins  []string      `env:"ALLOWED_ORIGINS" envSeparator:","`
}

func Load() Config {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Sprintf("config: %v", err))
	}
	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{"http://localhost:5173"}
	}
	return cfg
}

func (c Config) IsDevelopment() bool {
	return c.AppEnv == "development" || c.AppEnv == "dev"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
