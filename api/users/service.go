package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

type service struct {
}

var app *config.AppConfig

func New(a *config.AppConfig) service {
	app = a
	return service{}
}

// List retrieves all users
func (s *service) List(c echo.Context) error {
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
func (s *service) Get(c echo.Context) error {
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
func (s *service) Create(c echo.Context) error {

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
func (s *service) Delete(c echo.Context) error {

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
func (s *service) Update(c echo.Context) error {

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
	_, err := model.Update(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying user data",
			Value:   err.Error(),
		})
	}

	u, err := model.Get(login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error retrieving modified user",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: u,
	})
}
