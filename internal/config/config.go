package config

import (
	"time"

	"github.com/LikhithMar14/gopher-chat/pkg/env"
)

type Config struct {
	Addr        string
	DB          DBConfig
	Env         string
	APIURL      string
	FrontendURL string
	Mail        MailConfig
	FromEmail   string
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  string
}

type MailConfig struct {
	Sendgrid  SendgridConfig
	MailTrap  MailTrapConfig
	Exp       time.Duration
	FromEmail string
}

type SendgridConfig struct {
	APIKey string
}

type MailTrapConfig struct {
	APIKey string
}

func Load() Config {
	cfg := Config{
		Addr: env.GetString("PORT", ":8080"),
		DB: DBConfig{
			Addr:         env.GetString("DB_ADDR", ""),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  env.GetString("DB_MAX_LIFETIME", "15m"),
		},
		Env:         env.GetString("ENV", "development"),
		APIURL:      env.GetString("API_URL", "http://localhost:8080"),
		FrontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		FromEmail:   env.GetString("FROM_EMAIL", ""),
		Mail: MailConfig{
			Sendgrid: SendgridConfig{
				APIKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			MailTrap: MailTrapConfig{
				APIKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
			Exp: env.GetDuration("MAIL_EXP", 10*time.Minute),
		},
	}



	return cfg
}
