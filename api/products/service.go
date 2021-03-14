package products

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rfranzoia/cupuama-go/config"
	"github.com/rfranzoia/cupuama-go/utils"
)

type service struct {
}

var app *config.AppConfig

func NewProductService(a *config.AppConfig) service {
	app = a
	return service{}
}

// List retrieves all products
func (s *service) List(c echo.Context) error {
	list, err := model.List()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("Error searching Products"),
			Value:   err.Error(),
		})
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})
}

// Get retrieves an product by id
func (s *service) Get(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	u, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: fmt.Sprintf("Error searching Product %d", id),
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: u,
	})
}

// Create add a new product
func (s *service) Create(c echo.Context) error {

	product := new(Products)

	if err := c.Bind(product); err != nil {
		log.Println("(CreateProduct:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating product",
			Value:   err.Error(),
		})
	}

	id, err := model.Create(product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating product",
			Value:   err.Error(),
		})
	}

	u, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating product",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: u,
	})
}

// Delete removes an product by id
func (s *service) Delete(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := model.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error removing product",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "product successfully Deleted",
	})
}

// Update changes the data of an product
func (s *service) Update(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product := new(Products)

	if err := c.Bind(product); err != nil {
		log.Println("(Update:Bind)", err)
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying product data",
			Value:   err.Error(),
		})
	}

	product.ID = id
	_, err := model.Update(product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying product data",
			Value:   err.Error(),
		})
	}

	f, err := model.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error retrieving modified product",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.MessageJSON{
		Value: f,
	})
}
