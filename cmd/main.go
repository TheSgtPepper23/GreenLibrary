package main

import (
	"log"

	"github.com/TheSgtPepper23/GreenLibrary/db"
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

	// server.Renderer = models.NewTemplate()

	// server.Static("/static", "static")

	// ctx := models.NewContext()

	conn, err := db.GetConnection()
	if err != nil {
		server.StdLogger.Fatal()
	}

	defer conn.Close()

	sqlCollection := db.NewSQLCollectionContext(conn)

	server.GET("/collections", func(c echo.Context) error {
		collections, err := sqlCollection.SQLRetrieveCollections()
		if err != nil {
			server.Logger.Print(err.Error())
			return c.JSON(400, nil)
		}
		return c.JSON(200, collections)
	})

	// server.GET("/", func(c echo.Context) error {
	// 	collections, err := sqlCollection.SQLRetrieveCollections()
	// 	if err != nil {
	// 		server.Logger.Print(err.Error())
	// 		return c.Render(400, "index", 1)
	// 	}

	// 	ctx.Data = models.MainResponse{
	// 		Collections: collections,
	// 		Books:       &[]models.Book{},
	// 	}
	// 	return c.Render(200, "index", ctx)
	// })

	// server.POST("/collection", func(c echo.Context) error {
	// 	temp := models.Collection{
	// 		Name:         c.FormValue("name"),
	// 		CreationDate: time.Now(),
	// 	}
	// 	sqlCollection.SQLCreateCollection(&temp)

	// 	return c.Render(200, "collection-oob", temp)
	// })

	// server.GET("/book", func(c echo.Context) error {
	// 	books, err := services.SearchBook(c.FormValue("title"))
	// 	if err != nil {
	// 		return c.Render(400, "index", 1)
	// 	}
	// 	ctx.Data.ChangeBooks(books)
	// 	return c.Render(200, "books", ctx)
	// })

	server.Logger.Fatal(server.Start(":5555"))
}
