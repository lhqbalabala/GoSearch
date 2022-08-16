package controller

import (
	"RelatedSearch_Api/models"
	"RelatedSearch_Api/pkg/logger"
	"RelatedSearch_Api/pkg/weberror"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRelatedSearch(c *gin.Context) {
	query := c.Param("text")
	Service := c.Keys["prodservice"].(models.RelatedSearchService)
	text, err := Service.GetRelatedSearch(context.Background(), &models.RelatedSearchRequest{Query: query})
	if err != nil {
		logger.Logger.Errorf("[SayHello] failed to get text, err: %s", err.Error())
		c.JSON(http.StatusInternalServerError, weberror.Info{Error: http.StatusText(http.StatusBadRequest)})
		return
	}
	c.JSON(http.StatusOK, text)
}
