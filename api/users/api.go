package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

var app *config.AppConfig

func NewAPI(a *config.AppConfig) {
	app = a
}

// NewOrderAPI setups the configuration for orders
func RegisterRouting(g *echo.Group, a *config.AppConfig) {
	NewAPI(a)

	gu := g.Group("/v2/users")
	gu.GET("", List)
	gu.GET("/:login", Get)
	gu.POST("", Create)
	gu.PUT("/:login", Update)
	gu.DELETE("/:login", Delete)
}

// List retrieves all users
func List(c echo.Context) error {
	list, err := model.List()
	if err != nil {
		return nil
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})
}

// Get retrieves an user by login
func Get(c echo.Context) error {
	login := c.Param("login")

	u, err := model.Get(login)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("Error searching User %s", login),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: u,
	})
}

// Create add a new user
func Create(c echo.Context) error {

	user := new(Users)

	if err := c.Bind(user); err != nil {
		log.Println("(CreateUser:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating user",
			Value:   err.Error(),
		})
	}

	if err := model.Create(user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating user",
			Value:   err.Error(),
		})
	}

	login := user.Login
	u, err := model.Get(login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating user",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: u,
	})
}

// Delete removes an user by login
func Delete(c echo.Context) error {

	login := c.Param("login")

	if err := model.Delete(login); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error removing user",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "user successfully Deleted",
	})
}

// Update changes the data of an user
func Update(c echo.Context) error {

	login := c.Param("login")
	user := new(Users)

	if err := c.Bind(user); err != nil {
		log.Println("(Update:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying user data",
			Value:   err.Error(),
		})
	}

	user.Login = login
	u, err := model.Update(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying user data",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: u,
	})
}
