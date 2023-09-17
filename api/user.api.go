package api

import (
	"cupuama-go/config"
	"cupuama-go/service"

	"github.com/labstack/echo"
)

type UserApi struct {
	service service.UserService
}

func NewUserAPI(app *config.AppConfig) *UserApi {
	return &UserApi{service: service.NewUserService(app)}
}

func (api *UserApi) RegisterRouting(g *echo.Group) {

	gru := g.Group("/v2/users")
	gru.GET("", api.service.List)
	gru.GET("/:login", api.service.GetByLogin)

	gru.POST("", api.service.Create)

	gru.PUT("/:login", api.service.UpdateByLogin)
	gru.DELETE("/:login", api.service.DeleteByLogin)

	grl := g.Group("/v2/login")
	grl.POST("", api.service.Login)

}
