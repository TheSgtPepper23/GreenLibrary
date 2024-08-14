package main

import (
	"fmt"
	"log"
	"os"

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

	conn, err := db.GetConnection()
	if err != nil {
		server.StdLogger.Fatal()
	}

	defer conn.Close()

	collDB := db.NewSQLCollectionContext(conn)
	bookDB := db.NewSQLBookContext(conn)

	collServices := server.Group("/collection")
	bookServices := server.Group("/book")
	secret := os.Getenv("SECRET")
	// collServices.Use(echojwt.JWT([]byte(secret)))
	// bookServices.Use(echojwt.JWT([]byte(secret)))

	collServices.POST("", func(c echo.Context) error {
		data := new(models.Collection)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		data.ID, err = collDB.CreateCollection(data)

		if err != nil {
			return c.JSON(422, nil)
		}

		return c.JSON(200, data)
	})

	collServices.PUT("", func(c echo.Context) error {
		data := new(models.Collection)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		err := collDB.UpdateCollection(data)

		if err != nil {
			return c.JSON(400, nil)
		}

		return c.JSON(200, nil)
	})

	collServices.GET("", func(c echo.Context) error {
		collections, err := collDB.GetCollections()
		fmt.Println(secret)
		if err != nil {
			server.Logger.Print(err.Error())
			return c.JSON(400, nil)
		}
		return c.JSON(200, collections)
	})

	bookServices.POST("", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		data.ID, err = bookDB.CreateNewBook(data)
		if err != nil {
			return c.JSON(422, nil)
		}
		return c.JSON(200, data)
	})
	server.Logger.Fatal(server.Start(":5555"))
}
