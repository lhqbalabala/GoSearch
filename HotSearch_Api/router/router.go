package router

import (
	"HotSearch_Api/controller"
	middlewares "HotSearch_Api/middleware"
	"HotSearch_Api/models"
	"github.com/gin-gonic/gin"
)

func InitRouter(prodService models.HotSearchService) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.InitMiddleware(prodService), middlewares.ProdMiddleware(), middlewares.ErrorMiddleware())
	r.Handle("POST", "HotSearch/GetHotSearch", controller.GetHotSearch)
	r.Handle("POST", "HotSearch/InsertQuery", controller.InsertQuery)
	return r
}
