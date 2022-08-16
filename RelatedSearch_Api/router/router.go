package router

import (
	"RelatedSearch_Api/controller"
	middlewares "RelatedSearch_Api/middleware"
	"RelatedSearch_Api/models"
	"github.com/gin-gonic/gin"
)

func InitRouter(prodService models.RelatedSearchService) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.InitMiddleware(prodService), middlewares.ProdMiddleware(), middlewares.ErrorMiddleware())
	r.POST("/RelatedSearch/:text", controller.GetRelatedSearch) // 相关搜索
	return r
}
