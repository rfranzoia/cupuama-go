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

type FruitService struct {
	app        *config.AppConfig
	repository repository.FruitRepository
}

func NewFruitService(a *config.AppConfig) FruitService {
	return FruitService{
		app:        a,
		repository: repository.NewFruitRepository(a),
	}
}

// List retrieves all fruits
func (fs *FruitService) List(c echo.Context) error {
	list, err := fs.repository.List()
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.MessageJSON{
			Message: "Error searching Fruits",
			Value:   err.Error(),
		})
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})
}

// Get retrieves an fruit by id
func (fs *FruitService) Get(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	u, err := fs.repository.Get(id)
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
func (fs *FruitService) Create(c echo.Context) error {

	fruit := new(domain.Fruits)

	if err := c.Bind(fruit); err != nil {
		logger.Log.Info("(CreateFruit:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating fruit",
			Value:   err.Error(),
		})
	}

	id, err := fs.repository.Create(fruit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating fruit",
			Value:   err.Error(),
		})
	}

	u, err := fs.repository.Get(id)
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
func (fs *FruitService) Delete(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := fs.repository.Delete(id); err != nil {
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
func (fs *FruitService) Update(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	fruit := new(domain.Fruits)

	if err := c.Bind(fruit); err != nil {
		logger.Log.Info("(Update:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying fruit data",
			Value:   err.Error(),
		})
	}
	fruit.ID = id
	fmt.Println("parsed fruit", fruit)
	_, err := fs.repository.Update(fruit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying fruit data",
			Value:   err.Error(),
		})
	}

	f, err := fs.repository.Get(id)
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
