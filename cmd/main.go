package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TheSgtPepper23/GreenLibrary/db"
	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file j")
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

	secret := os.Getenv("SECRET")
	collServices := server.Group("/collection", echojwt.JWT([]byte(secret)))
	bookServices := server.Group("/book", echojwt.JWT([]byte(secret)))
	authServices := server.Group("/auth")
	adminServices := server.Group("/admin", echojwt.JWT([]byte(secret)))

	server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	collServices.POST("", func(c echo.Context) error {
		data := new(models.Collection)
		if err := c.Bind(data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}

		err = collDB.CreateCollection(data)

		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
		}

		return c.JSON(200, data)
	})

	collServices.PUT("", func(c echo.Context) error {
		data := new(models.Collection)
		if err := c.Bind(data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}

		err := collDB.UpdateCollection(data)

		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
		}

		return c.JSON(200, data)
	})

	collServices.GET("/:userID", func(c echo.Context) error {
		stringID := c.Param("userID")
		collections, err := collDB.GetCollections(stringID)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrNotFound
		}
		return c.JSON(200, collections)
	})

	collServices.DELETE("/:collection", func(c echo.Context) error {
		stringID := c.Param("collection")
		err := collDB.DeleteCollection(stringID)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrNotFound
		}
		return c.JSON(200, nil)
	})

	bookServices.POST("", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		fmt.Println(claims)

		err = bookDB.CreateNewBook(data, claims["userKey"].(string))
		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
		}
		return c.JSON(200, data)
	})

	bookServices.PUT("", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}

		err = bookDB.UpdateBook(data)
		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
		}
		return c.JSON(200, data)
	})

	bookServices.GET("/:collection", func(c echo.Context) error {
		stringID := c.Param("collection")

		books, err := bookDB.GetBooksOfCollection(stringID)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrNotFound
		}
		return c.JSON(200, books)
	})

	bookServices.POST("/search", func(c echo.Context) error {
		data := make(map[string]string)
		if err := c.Bind(&data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}
		results, err := services.SearchBook(data["title"])
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrServiceUnavailable
		}
		return c.JSON(200, results)
	})

	bookServices.PUT("/delete", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}

		err = bookDB.RemoveBookFromCollection(data)
		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
		}
		return c.JSON(200, data)
	})

	bookServices.PUT("/move", func(c echo.Context) error {
		data := new(models.Book)
		if err := c.Bind(data); err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}

		err = bookDB.MoveBook(data)
		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
		}

		return c.JSON(200, data)
	})

	authServices.POST("/login", func(c echo.Context) error {
		userData := new(models.User)
		err := c.Bind(userData)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}
		userKey, err := userDB.AuthenticateUser(userData)

		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrUnauthorized
		}

		token, err := services.GenerateToken(userData.Email, userKey)

		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrUnauthorized
		}

		return c.JSON(200, token)
	})

	authServices.POST("/refresh", func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		var newToken string
		if token != "" {
			newToken, err = services.RefreshToken(token)
			if err != nil {
				fmt.Println(err.Error())
				return echo.ErrUnauthorized
			}
		}
		return c.JSON(200, newToken)
	})

	adminServices.POST("/register", func(c echo.Context) error {
		userData := new(models.User)
		err := c.Bind(userData)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrBadRequest
		}
		err = userDB.UserWizard(userData.Email, userData.Password)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrUnprocessableEntity
		}
		return c.JSON(200, userData)
	})

	server.Logger.Fatal(server.Start(":5555"))
}
