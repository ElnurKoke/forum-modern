package main

import (
	"forum/internal/handler"
	"forum/internal/models"
	svr "forum/internal/server"
	"forum/internal/service"
	"forum/internal/storage"
)

func main() {
	config, err := svr.NewConfig()
	if err != nil {
		models.ErrLog.Println(err)
	}
	db := storage.InitDB(config)
	store := storage.NewStorage(db)
	services := service.NewService(store)
	handlers := handler.NewHandler(services)
	server := new(svr.Server)
	if err := server.Run(config.Port, handlers.InitRoutes()); err != nil {
		models.ErrLog.Printf("Error running server: %s\n", err)
	}
}
