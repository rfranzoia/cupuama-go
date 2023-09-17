package api

import (
	"cupuama-go/config"
	"cupuama-go/service"
	"github.com/labstack/echo"
)

type ProductApi struct {
	service service.ProductService
}

func NewProductAPI(app *config.AppConfig) *ProductApi {
	return &ProductApi{service: service.NewProductService(app)}
}

func (api *ProductApi) RegisterRouting(g *echo.Group) {

	grp := g.Group("/v2/products")
	grp.GET("", api.service.List)
	grp.GET("/:id", api.service.Get)
	grp.POST("", api.service.Create)
	grp.PUT("/:id", api.service.Update)
	grp.DELETE("/:id", api.service.Delete)

}
