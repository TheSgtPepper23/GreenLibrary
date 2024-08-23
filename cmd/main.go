package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TheSgtPepper23/GreenLibrary/db"
	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
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

	//TODO quitar las líneas de abajo, a la hora de desplegarlo se usará un reverse proxy
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allowed origins
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true, // Set to true if your API requires credentials (e.g., cookies)
	}))

	conn, err := db.GetConnection()
	if err != nil {
		server.StdLogger.Fatal()
	}

	defer conn.Close()

	collDB := db.NewSQLCollectionContext(conn)
	bookDB := db.NewSQLBookContext(conn)
	userDB := db.NewSQLUserContext(conn)

	collServices := server.Group("/collection")
	bookServices := server.Group("/book")
	authServices := server.Group("/auth")
	secret := os.Getenv("SECRET")
	// collServices.Use(echojwt.JWT([]byte(secret)))
	// bookServices.Use(echojwt.JWT([]byte(secret)))

	collServices.POST("", func(c echo.Context) error {
		data := new(models.Collection)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		err = collDB.CreateCollection(data)

		if err != nil {
			fmt.Println(err.Error())
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

		return c.JSON(200, data)
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

		err = bookDB.CreateNewBook(data)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(422, nil)
		}
		return c.JSON(200, data)
	})

	bookServices.PUT("", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		err = bookDB.UpdateBook(data)
		if err != nil {
			return c.JSON(422, nil)
		}
		return c.JSON(200, data)
	})

	bookServices.GET("/:collection", func(c echo.Context) error {
		stringID := c.Param("collection")

		if err != nil {
			return c.JSON(400, nil)
		}

		books, err := bookDB.GetBookOfCollection(stringID)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(500, nil)
		}
		return c.JSON(200, books)
	})

	bookServices.POST("/search", func(c echo.Context) error {
		data := make(map[string]string)
		if err := c.Bind(&data); err != nil {
			return c.JSON(400, nil)
		}
		results, err := services.SearchBook(data["title"])
		if err != nil {
			return c.JSON(400, nil)
		}
		return c.JSON(200, results)
	})

	authServices.POST("/login", func(c echo.Context) error {
		userData := new(models.User)
		err := c.Bind(userData)
		if err != nil {
			return c.JSON(400, nil)
		}
		err = userDB.AuthenticateUser(userData)

		if err != nil {
			return c.JSON(401, nil)
		}

		token, err := services.GenerateToken(userData.Email)

		if err != nil {
			return c.JSON(401, nil)
		}

		return c.JSON(200, token)
	})

	authServices.POST("/refresh", func(c echo.Context) error {
		return c.JSON(200, nil)
	})

	server.Logger.Fatal(server.Start(":5555"))
}

//TODO probar las funciones de SQL
//TODO Agregar un struct de respuesta estandarizada a todos los endpoints
//TODO habilitar nuevamente la autenticación
