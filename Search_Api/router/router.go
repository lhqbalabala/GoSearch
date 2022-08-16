package router

import (
	"Search_Api/controller"
	middlewares "Search_Api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ProdMiddleware(), middlewares.ErrorMiddleware())
	r.Handle("POST", "/query", controller.Search)
	r.Handle("POST", "/query/picture", controller.SearchPicture)
	return r
}
