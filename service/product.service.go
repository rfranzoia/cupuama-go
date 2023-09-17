package service

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"cupuama-go/repository"
	"cupuama-go/utils"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type ProductService struct {
	app        *config.AppConfig
	repository repository.ProductRepository
}

func NewProductService(a *config.AppConfig) ProductService {
	return ProductService{
		app:        a,
		repository: repository.NewProductRepository(a),
	}
}

// List retrieves all products
func (ps *ProductService) List(c echo.Context) error {
	list, err := ps.repository.List()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: "Error searching Products",
			Value:   err.Error(),
		})
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})
}

// Get retrieves an product by id
func (ps *ProductService) Get(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	u, err := ps.repository.Get(id)
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
func (ps *ProductService) Create(c echo.Context) error {

	product := new(domain.Products)

	if err := c.Bind(product); err != nil {
		logger.Log.Info("(CreateProduct:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating product",
			Value:   err.Error(),
		})
	}

	id, err := ps.repository.Create(product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating product",
			Value:   err.Error(),
		})
	}

	u, err := ps.repository.Get(id)
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
func (ps *ProductService) Delete(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := ps.repository.Delete(id); err != nil {
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
func (ps *ProductService) Update(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product := new(domain.Products)

	if err := c.Bind(product); err != nil {
		logger.Log.Info("(Update:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying product data",
			Value:   err.Error(),
		})
	}

	product.ID = id
	_, err := ps.repository.Update(product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying product data",
			Value:   err.Error(),
		})
	}

	f, err := ps.repository.Get(id)
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
