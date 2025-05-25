package main

import (
	"log"

	 "github.com/LikhithMar14/gopher-chat/cmd/migration"
	"github.com/LikhithMar14/gopher-chat/internal/db"
	"github.com/LikhithMar14/gopher-chat/internal/env"
	"github.com/LikhithMar14/gopher-chat/internal/store"
)





func main() {

	cfg := config{
		addr: env.GetString("PORT", ":8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5432/gopher-chat?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 50),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxLifetime: env.GetString("DB_MAX_LIFETIME", "15m"),
		},
		
	}
	database,err := db.Open(cfg.db.addr,cfg.db.maxOpenConns,cfg.db.maxIdleConns,cfg.db.maxLifetime)
	if err != nil {
		log.Fatal(err)
	}
	err = db.MigrateFS(database,migration.FS,".")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()
	log.Println("database connection pool established")

	

	store := store.NewStorage(database)

	app := &application{
		config: cfg,
		store: store,
	}


	mux := app.mount()

	log.Fatal(app.run(mux))

}
