package main

import (
	"log"

	"github.com/LikhithMar14/gopher-chat/internal/api"
	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/db"
	"github.com/LikhithMar14/gopher-chat/internal/migrations"
	"github.com/LikhithMar14/gopher-chat/internal/store"
)

func main() {

	cfg := config.Load()


	database, err := db.Open(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxLifetime)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()


	err = db.MigrateFS(database, migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connection pool established")

	storage := store.NewStorage(database)

	app := api.NewApplication(cfg, storage)

	mux := app.Routes()

	log.Fatal(app.Serve(mux))
}
