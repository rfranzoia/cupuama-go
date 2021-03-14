package fruits

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	gu := g.Group("/v2/fruits")
	gu.GET("", List)
	gu.GET("/:id", Get)
	gu.POST("", Create)
	gu.PUT("/:id", Update)
	gu.DELETE("/:id", Delete)
}

// List retrieves all fruits
func List(c echo.Context) error {
	list, err := model.List()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("Error searching Fruits"),
			Value:   err.Error(),
		})
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})
}

// Get retrieves an fruit by id
func Get(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	u, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("Error searching Fruit %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: u,
	})
}

// Create add a new fruit
func Create(c echo.Context) error {

	fruit := new(Fruits)

	if err := c.Bind(fruit); err != nil {
		log.Println("(CreateFruit:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating fruit",
			Value:   err.Error(),
		})
	}

	id, err := model.Create(fruit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating fruit",
			Value:   err.Error(),
		})
	}

	u, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating fruit",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: u,
	})
}

// Delete removes an fruit by id
func Delete(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := model.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error removing fruit",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "fruit successfully Deleted",
	})
}

// Update changes the data of an fruit
func Update(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fruit := new(Fruits)

	if err := c.Bind(fruit); err != nil {
		log.Println("(Update:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying fruit data",
			Value:   err.Error(),
		})
	}

	fruit.ID = id
	_, err := model.Update(fruit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying fruit data",
			Value:   err.Error(),
		})
	}

	f, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error retrieving modified fruit",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: f,
	})
}
