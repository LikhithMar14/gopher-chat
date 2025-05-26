package config

import "github.com/LikhithMar14/gopher-chat/internal/env"

type Config struct {
	Addr string
	DB   DBConfig
	Env  string
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  string
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
	}
}
