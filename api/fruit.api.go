package api

import (
	"cupuama-go/config"
	"cupuama-go/service"
	"github.com/labstack/echo"
)

type FruitApi struct {
	service service.FruitService
}

func NewFruitAPI(app *config.AppConfig) *FruitApi {
	return &FruitApi{service: service.NewFruitService(app)}
}

func (api *FruitApi) RegisterRouting(g *echo.Group) {

	grp := g.Group("/v2/fruits")
	grp.GET("", api.service.List)
	grp.GET("/:id", api.service.Get)
	grp.POST("", api.service.Create)
	grp.PUT("/:id", api.service.Update)
	grp.DELETE("/:id", api.service.Delete)

}
