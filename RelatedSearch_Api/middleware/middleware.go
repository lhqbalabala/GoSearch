package middlewares

import (
	"RelatedSearch_Api/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// InitMiddleware 注入prodService中间件
func InitMiddleware(prodService models.RelatedSearchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		c.Keys["prodservice"] = prodService
	}
}
func ProdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("mytoken") != "\"lgdSearch\"" {
			c.JSON(500, gin.H{"status": fmt.Sprintf("%s", "禁止外部访问!")})
			c.Abort()
		}
		c.Next()
	}
}

// ErrorMiddleware 异常处理
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(500, gin.H{"status": fmt.Sprintf("%s", r)})
				c.Abort()
			}
		}()
		c.Next()
	}
}
