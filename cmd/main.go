package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TheSgtPepper23/GreenLibrary/db"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DatabaseContext struct {
	CollDB *db.CollectionSQLContext
	BookDb *db.BookSQLContext
	UserDB *db.UserSQLContext
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file j")
	}
	secret := os.Getenv("SECRET")

	server := echo.New()

	conn, err := db.GetConnection()
	if err != nil {
		server.StdLogger.Fatal()
	}
	defer conn.Close()
	dbContext := &DatabaseContext{
		CollDB: db.NewSQLCollectionContext(conn),
		BookDb: db.NewSQLBookContext(conn),
		UserDB: db.NewSQLUserContext(conn),
	}

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "https://andresdglez.com"}, // Allowed origins
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true, // Set to true if your API requires credentials (e.g., cookies)
	}))
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("dbContext", dbContext)
			return next(c)
		}
	})

	server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	//Collection endpoints
	collServices := server.Group("/collection", echojwt.JWT([]byte(secret)))
	collServices.POST("", HandlerCreateCollection)
	collServices.PUT("", HandlerUpdateCollection)
	collServices.GET("/:userID", HandlerGetCollections)
	collServices.DELETE("/:collection", HandlerDeleteCollection)

	//Book endpoints
	bookServices := server.Group("/book", echojwt.JWT([]byte(secret)))
	bookServices.POST("", HandlerCreateNewBook)
	bookServices.PUT("", HandlerUpdateBook)
	bookServices.GET("/:collection", HandlerGetCollectonBooks)
	bookServices.POST("/search", HandlerSearchBook)
	bookServices.PUT("/delete", HandlerRemoveFromCollection)
	bookServices.PUT("/move", HandlerMoveBook)

	//Auth endpoints
	authServices := server.Group("/auth")
	authServices.POST("/login", HandlerLogin)
	authServices.POST("/refresh", HandlerRefreshToken)

	//Admin endpoints
	adminServices := server.Group("/admin", echojwt.JWT([]byte(secret)))
	adminServices.POST("/register", HandlerRegister)

	server.Logger.Fatal(server.Start(":5555"))
}
