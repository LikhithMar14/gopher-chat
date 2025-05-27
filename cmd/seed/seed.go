package main

import (
	"log"

	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/utils/env"
	"github.com/LikhithMar14/gopher-chat/internal/db"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/store"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxLifetime)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	
	log.Println("Seeding database...")

	if env.GetString("ENV", "development") == "development" {
		_,err = database.Exec(`
			TRUNCATE TABLE users CASCADE;
			TRUNCATE TABLE posts CASCADE;
			TRUNCATE TABLE comments CASCADE;
		`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Truncated tables")
	}

	

	storage := store.NewStorage(database)

		

	userService := service.NewUserService(storage)
	postService := service.NewPostService(storage)
	commentService := service.NewCommentService(storage)

	err = db.Seed(database, userService, postService, commentService)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database seeding completed successfully!")
}
