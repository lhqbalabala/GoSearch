package controller

import (
	"HotSearch_Api/handler"
	"HotSearch_Api/models"
	"HotSearch_Api/payloads"
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHotSearch(c *gin.Context) {
	flow_Entry, flow_err := sentinel.Entry("GetHotSearch_flow", sentinel.WithTrafficType(base.Inbound))
	if flow_err != nil {
		c.JSON(http.StatusOK, payloads.Success(handler.GetHotSearchCache()))
		return
	}
	flow_Entry.Exit()
	Service := c.Keys["prodservice"].(models.HotSearchService)
	var Req models.HotSearchRequest
	err := c.Bind(&Req)
	if err != nil {
		c.JSON(500, gin.H{
			"status": err.Error()})
		return
	}
	break_Entry, break_err := sentinel.Entry("GetHotSearch_break", sentinel.WithTrafficType(base.Inbound))
	if break_err != nil {
		c.JSON(http.StatusOK, payloads.Success(handler.GetHotSearchCache()))
		return
	}
	Res, err := Service.GetHotSearchList(context.Background(), &models.HotSearchRequest{Size: 10})
	if err != nil {
		c.JSON(500, gin.H{
			"status": err.Error()})
		return
	}
	handler.SetHotSearchCache(Res)
	c.JSON(http.StatusOK, payloads.Success(Res))
	break_Entry.Exit()

}
