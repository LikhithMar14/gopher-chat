package config

import (
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/utils/env"
)

type Config struct {
	Addr string
	DB   DBConfig
	Env  string
	APIURL string
	Mail   MailConfig
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  string
}

type MailConfig struct {
	Exp time.Duration
}
func Load() Config {
	return Config{
		Addr: env.GetString("PORT", ":8080"),
		DB: DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5432/gopher-chat?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  env.GetString("DB_MAX_LIFETIME", "15m"),
		},
		Env: env.GetString("ENV", "development"),
		APIURL: env.GetString("API_URL", "http://localhost:8080"),
		Mail: MailConfig{
			Exp: env.GetDuration("MAIL_EXP", 10*time.Minute),
		},
	}
}
