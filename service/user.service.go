package service

import (
	"cupuama-go/config"
	"cupuama-go/domain"
	"cupuama-go/logger"
	"cupuama-go/repository"
	"cupuama-go/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type UserService struct {
	app        *config.AppConfig
	repository repository.UserRepository
}

func NewUserService(a *config.AppConfig) UserService {
	return UserService{
		app:        a,
		repository: repository.NewUserRepository(a),
	}
}

// Login validates an user
func (s *UserService) Login(c echo.Context) error {

	user := new(domain.Users)

	if err := c.Bind(user); err != nil {
		logger.Log.Info("(Login:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "invalid username/password",
			Value:   err.Error(),
		})
	}

	logger.Log.Info("user: " + user.Login)
	u, err := s.repository.GetByLogin(user.Login)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.MessageJSON{
			Message: "invalid username/password",
			Value:   err.Error(),
		})
	}
	logger.Log.Info("u: " + u.Login)
	if u.Login != user.Login || u.Password != user.Password {
		return c.JSON(http.StatusUnauthorized, utils.MessageJSON{
			Message: "invalid username/password",
			Value:   nil,
		})
	}

	token, err := utils.CreateJwtToken(user.Login, user.Person.FirstName)
	if err != nil {
		logger.Log.Info("Erro ao criar Token JWT")
		return c.JSON(http.StatusInternalServerError, utils.MessageJSON{
			Message: "Erro ao criar Token JWT",
			Value:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Message: "login successful!",
		Value:   token,
	})
}

// List retrieves all users
func (s *UserService) List(c echo.Context) error {
	list, err := s.repository.FindAll()
	if err != nil {
		return nil
	}
	defer c.Request().Body.Close()

	return c.JSON(http.StatusOK, utils.MessageJSON{
		Value: list,
	})
}

// Get retrieves an user by login
func (s *UserService) GetByLogin(c echo.Context) error {
	login := c.Param("login")

	u, err := s.repository.GetByLogin(login)
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
func (s *UserService) Create(c echo.Context) error {

	user := new(domain.Users)

	if err := c.Bind(user); err != nil {
		logger.Log.Info("(CreateUser:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating user",
			Value:   err.Error(),
		})
	}

	if err := s.repository.Create(user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error creating user",
			Value:   err.Error(),
		})
	}

	login := user.Login
	u, err := s.repository.GetByLogin(login)
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
func (s *UserService) DeleteByLogin(c echo.Context) error {

	login := c.Param("login")

	if err := s.repository.DeleteByLogin(login); err != nil {
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
func (s *UserService) UpdateByLogin(c echo.Context) error {

	login := c.Param("login")
	user := new(domain.Users)

	if err := c.Bind(user); err != nil {
		logger.Log.Info("(Update:Bind)" + err.Error())
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying user data",
			Value:   err.Error(),
		})
	}

	user.Login = login
	_, err := s.repository.UpdateByLogin(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.MessageJSON{
			Message: "Error modifying user data",
			Value:   err.Error(),
		})
	}

	u, err := s.repository.GetByLogin(login)
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
