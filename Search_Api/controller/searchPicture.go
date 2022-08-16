package controller

import (
	"Search_Api/handler"
	"Search_Api/models"
	"Search_Api/payloads"
	"Search_Api/pkg/grpcExpand"
	balancer "Search_Api/pkg/grpcExpand/balance"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchPicture(c *gin.Context) {
	var Req models.SearchRequest
	err := c.Bind(&Req)
	if err != nil {
		c.JSON(500, gin.H{
			"status": err.Error()})
		return
	}
	Req.Likes = handler.GetAllFavoritesDocs(c)
	handler.InsertQuery(Req.Query)
	rsp, err := grpcExpand.Client().GetSearchPictureResult(context.WithValue(context.Background(), balancer.ConsistentHash, Req.Query), &Req)
	//handler.AddRelatedSearch(Req.Query, rsp)

	c.JSON(http.StatusOK, payloads.Success(rsp))
}
