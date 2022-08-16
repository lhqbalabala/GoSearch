package controller

import (
	"HotSearch_Api/models"
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertQuery(c *gin.Context) {
	flow_Entry, flow_err := sentinel.Entry("InsertQuery_flow", sentinel.WithTrafficType(base.Inbound))
	if flow_err != nil {
		c.JSON(500, "限流了")
		return
	}
	flow_Entry.Exit()
	Service := c.Keys["prodservice"].(models.HotSearchService)
	var Req models.InsertQueryRequest
	err := c.Bind(&Req)
	if err != nil {
		c.JSON(500, gin.H{
			"status": err.Error()})
		return
	}
	break_Entry, break_err := sentinel.Entry("GetHotSearch_break", sentinel.WithTrafficType(base.Inbound))
	if break_err != nil {
		c.JSON(500, "熔断了")
		return
	}
	_, err = Service.InsertQuery(context.Background(), &models.InsertQueryRequest{Query: Req.Query})
	if err != nil {
		c.JSON(500, gin.H{
			"status": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
	break_Entry.Exit()

}
