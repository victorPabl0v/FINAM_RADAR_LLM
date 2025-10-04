package main

import (
	"FINAM_RADAR/db"
	"FINAM_RADAR/internal/handlers"
	"FINAM_RADAR/internal/repository"
	"FINAM_RADAR/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	pool := db.NewPool()
	defer pool.Close()

	e := echo.New()

	newsHandler := handlers.NewNewsHandler(
		service.NewNewsService(repository.NewNewsRepository(pool)),
	)

	newsGroup := e.Group("/news")
	newsGroup.POST("", newsHandler.AddNews)
	newsGroup.GET("", newsHandler.GetNews)

	e.Logger.Fatal(e.Start(":8080"))
}
