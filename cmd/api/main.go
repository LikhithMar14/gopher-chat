package main

import (
	"github.com/LikhithMar14/gopher-chat/internal/api"
	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/db"
	"github.com/LikhithMar14/gopher-chat/internal/migrations"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"go.uber.org/zap"
)

const Version = "0.0.1"

//	@title			Gopher Chat API
//	@version		1.0.0
//	@description	This is a sample server Gopher Chat API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/v1

//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for the API Key

func main() {
	//http://localhost:8080/v1/swagger/index.html
	cfg := config.Load()

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	database, err := db.Open(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxLifetime)
	if err != nil {
		logger.Fatalw("Failed to open database connection", "error", err)
	}
	defer database.Close()

	// should give the path related to where you are doing goose.Up
	err = db.MigrateFS(database, migrations.FS, ".")
	if err != nil {
		logger.Fatalw("Failed to migrate database", "error", err)
	}

	logger.Infow("Database connection pool established and migrated successfully", "addr", cfg.DB.Addr, "maxOpenConns", cfg.DB.MaxOpenConns, "maxIdleConns", cfg.DB.MaxIdleConns, "maxLifetime", cfg.DB.MaxLifetime)

	storage := store.NewStorage(database)

	app := api.NewApplication(cfg, storage, Version, logger)

	mux := app.Routes()

	logger.Fatal(app.Serve(mux))
}
