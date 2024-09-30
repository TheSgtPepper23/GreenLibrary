package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func HandlerCreateCollection(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := new(models.Collection)
	if err := c.Bind(data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	err := dbContext.CollDB.CreateCollection(data)

	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
	}

	return c.JSON(200, data)
}

func HandlerUpdateCollection(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := new(models.Collection)
	if err := c.Bind(data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	err := dbContext.CollDB.UpdateCollection(data)

	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
	}

	return c.JSON(200, data)
}

func HandlerGetCollections(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	stringID := c.Param("userID")
	collections, err := dbContext.CollDB.GetCollections(stringID)
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrNotFound
	}
	return c.JSON(200, collections)
}

func HandlerDeleteCollection(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	stringID := c.Param("collection")
	err := dbContext.CollDB.DeleteCollection(stringID)
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrNotFound
	}
	return c.JSON(200, nil)
}

func HandlerCreateNewBook(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := new(models.Book)
	if err := c.Bind(data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	err := dbContext.BookDb.CreateNewBook(data, claims["userKey"].(string))
	if err != nil {
		if err.Error() == "book already read" {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "El libro ya está marcado como leído")
		}
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
	}
	return c.JSON(200, data)
}

func HandlerUpdateBook(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := new(models.Book)
	if err := c.Bind(data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	err := dbContext.BookDb.UpdateBook(data)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
	}
	return c.JSON(200, data)
}

// tiene los params ammount, page, y order
func HandlerGetCollectonBooks(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	stringID := c.Param("collection")
	ammout := c.QueryParam("ammount")
	page := c.QueryParam("page")
	order := c.QueryParam("order")

	results, err := services.StringsToInts(ammout, page, order)
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	books, err := dbContext.BookDb.GetBooksOfCollection(stringID, results[0], results[1], models.OrderOption(results[2]))
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrNotFound
	}
	return c.JSON(200, books)
}

func HandlerSearchBook(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := make(map[string]string)
	if err := c.Bind(&data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	results := make([]models.Book, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 2)

	start := time.Now()
	wg.Add(2)
	go services.SearchBook(data["title"], &results, &wg, &mu, errChan)
	go dbContext.BookDb.SearchBookLocally(data["title"], &results, &wg, &mu, errChan)

	wg.Wait()
	close(errChan)
	fmt.Printf("total : %v", time.Since(start).Milliseconds())

	hasError := false
	for err := range errChan {
		if err != nil {
			fmt.Println(err.Error())
			hasError = true
			break
		}
	}

	if hasError && len(results) == 0 {
		return echo.ErrServiceUnavailable
	}

	return c.JSON(200, results)
}

func HandlerRemoveFromCollection(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := new(models.Book)
	if err := c.Bind(data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	err := dbContext.BookDb.RemoveBookFromCollection(data)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
	}
	return c.JSON(200, data)
}

func HandlerMoveBook(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	data := new(models.Book)
	if err := c.Bind(data); err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	err := dbContext.BookDb.MoveBook(data, claims["userKey"].(string))
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "No es posible realizar la operación")
	}

	return c.JSON(200, data)
}

func HandlerLogin(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	userData := new(models.User)
	err := c.Bind(userData)
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}
	userKey, err := dbContext.UserDB.AuthenticateUser(userData)

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
}

func HandlerRefreshToken(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token != "" {
		newToken, err := services.RefreshToken(token)
		if err != nil {
			fmt.Println(err.Error())
			return echo.ErrUnauthorized
		}
		return c.JSON(200, newToken)
	}
	return echo.ErrBadRequest
}

func HandlerRegister(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)
	userData := new(models.User)
	err := c.Bind(userData)
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrBadRequest
	}
	err = dbContext.UserDB.UserWizard(userData.Email, userData.Password)
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrUnprocessableEntity
	}
	return c.JSON(200, userData)
}

func HandlerGetLibrary(c echo.Context) error {
	dbContext := c.Get("dbContext").(*DatabaseContext)

	library, err := dbContext.BookDb.GetAllStoredBooks()
	if err != nil {
		fmt.Println(err.Error())
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, library)

}
