package main

import (
	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	db "github.com/LikhithMar14/gopher-chat/internal/store/database"
	"github.com/LikhithMar14/gopher-chat/internal/utils/env"
	mailer "github.com/LikhithMar14/gopher-chat/internal/utils/mailer"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	cfg := config.Load()

	database, err := db.Open(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxLifetime)
	if err != nil {
		logger.Fatal("Failed to open database connection", zap.Error(err))
	}
	defer database.Close()

	logger.Info("Seeding database...")

	if env.GetString("ENV", "development") == "development" {
		_, err = database.Exec(`
			TRUNCATE TABLE users CASCADE;
			TRUNCATE TABLE posts CASCADE;
			TRUNCATE TABLE comments CASCADE;
		`)
		if err != nil {
			logger.Fatal("Failed to truncate tables", zap.Error(err))
		}
		logger.Info("Truncated tables")
	}

	storage := store.NewStorage(database)

	postService := service.NewPostService(storage)
	commentService := service.NewCommentService(storage)
	mailer := mailer.NewSendgrid(cfg.Mail.Sendgrid.APIKey, cfg.FromEmail)
	authService := service.NewAuthService(storage, cfg.Mail.Exp, mailer, cfg, logger)

	err = db.Seed(database, authService, postService, commentService, logger)
	if err != nil {
		logger.Fatal("Failed to seed database", zap.Error(err))
	}

	logger.Info("Database seeding completed successfully!")
}
