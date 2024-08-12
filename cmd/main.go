package main

import (
	"log"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/db"
	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.Renderer = models.NewTemplate()

	server.Static("/static", "static")

	ctx := models.NewContext()

	conn, err := db.GetConnection()
	if err != nil {
		server.StdLogger.Fatal()
	}

	defer conn.Close()

	sqlCollection := db.NewSQLCollectionContext(conn)

	server.GET("/", func(c echo.Context) error {
		collections, err := sqlCollection.SQLRetrieveCollections()
		if err != nil {
			server.Logger.Print(err.Error())
			return c.Render(400, "index", 1)
		}

		resp := models.MainResponse{
			Collections: collections,
		}
		ctx.Data = resp
		return c.Render(200, "index", ctx)
	})

	server.POST("/collection", func(c echo.Context) error {
		temp := models.Collection{
			Name:         c.FormValue("name"),
			CreationDate: time.Now(),
		}
		sqlCollection.SQLCreateCollection(&temp)

		return c.Render(200, "collection-oob", temp)
	})

	server.Logger.Fatal(server.Start(":5555"))
}
