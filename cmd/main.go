package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TheSgtPepper23/GreenLibrary/db"
	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
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

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "https://andresdglez.com"}, // Allowed origins
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
	collServices.Use(echojwt.JWT([]byte(secret)))
	bookServices.Use(echojwt.JWT([]byte(secret)))

	server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	collServices.POST("", func(c echo.Context) error {
		data := new(models.Collection)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		err = collDB.CreateCollection(data)

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
			return c.JSON(422, nil)
		}

		return c.JSON(200, data)
	})

	collServices.GET("", func(c echo.Context) error {
		collections, err := collDB.GetCollections()
		if err != nil {
			return c.JSON(404, nil)
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

		books, err := bookDB.GetBooksOfCollection(stringID)
		if err != nil {
			return c.JSON(404, nil)
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
			return c.JSON(503, nil)
		}
		return c.JSON(200, results)
	})

	bookServices.PUT("/delete", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			return c.JSON(400, nil)
		}

		err = bookDB.RemoveBookFromCollection(data)
		if err != nil {
			return c.JSON(422, nil)
		}
		return c.JSON(200, data)
	})

	authServices.POST("/login", func(c echo.Context) error {
		userData := new(models.User)
		err := c.Bind(userData)
		if err != nil {
			return c.JSON(400, nil)
		}
		userKey, err := userDB.AuthenticateUser(userData)

		if err != nil {
			return c.JSON(401, nil)
		}

		token, err := services.GenerateToken(userData.Email, userKey)

		if err != nil {
			return c.JSON(401, nil)
		}

		return c.JSON(200, token)
	})

	authServices.POST("/refresh", func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		var newToken string
		if token != "" {
			newToken, err = services.RefreshToken(token)
			if err != nil {
				return c.JSON(401, nil)
			}
		}
		return c.JSON(200, newToken)
	})

	server.Logger.Fatal(server.Start(":5555"))
}
